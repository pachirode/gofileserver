package staticserver

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"

	goplist "github.com/fork2fix/go-plist"
	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/core"
)

type plistBundle struct {
	CFBundleIdentifier  string `plist:"CFBundleIdentifier"`
	CFBundleVersion     string `plist:"CFBundleVersion"`
	CFBundleDisplayName string `plist:"CFBundleDisplayName"`
	CFBundleName        string `plist:"CFBundleName"`
	CFBundleIconFile    string `plist:"CFBundleIconFile"`
	CFBundleIcons       struct {
		CFBundlePrimaryIcon struct {
			CFBundleIconFiles []string `plist:"CFBundleIconFiles"`
		} `plist:"CFBundlePrimaryIcon"`
	} `plist:"CFBundleIcons"`
}

type plAsset struct {
	Kind string `plist:"kind"`
	URL  string `plist:"url"`
}

type plItem struct {
	Assets   []*plAsset `plist:"assets"`
	Metadata struct {
		BundleIdentifier string `plist:"bundle-identifier"`
		BundleVersion    string `plist:"bundle-version"`
		Kind             string `plist:"kind"`
		Title            string `plist:"title"`
	} `plist:"metadata"`
}

type downloadPlist struct {
	Items []*plItem `plist:"items"`
}

func (s *HTTPStaticServer) PList(ctx *gin.Context) {
	requestPath := ctx.Request.URL.Path
	realPath := s.getRealPath(requestPath)

	if filepath.Ext(requestPath) == ".plist" {
		requestPath = requestPath[:len(requestPath)-6] + ".ipa"
	}

	plistInfo, err := parseIPA(realPath)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := &url.URL{
		Scheme: scheme,
		Host:   ctx.Request.Host,
	}
	data, err := generateDownloadPlist(baseURL, requestPath, plistInfo)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	ctx.XML(200, data)
}

func (s *HTTPStaticServer) IPALink(ctx *gin.Context) {
	requestPath := ctx.Request.URL.Path
	var plistURL string

	if ctx.Request.URL.Scheme == "https" {
		plistURL = combineURL(ctx, "/-/ipa/plist/"+requestPath).String()
	} else if s.PlistProxy != "" {
		httpPlistLink := "http://" + ctx.Request.Host + "/-/ipa/plist/" + requestPath
		url, err := s.genPlistLink(httpPlistLink)
		if err != nil {
			core.WriteResponse(ctx, err, nil)
			return
		}
		plistURL = url
	} else {
		core.WriteResponse(ctx, errors.New("500: Server should be https:// or provide valid plistproxy"), nil)
		return
	}

	renderHTML(ctx, "assets/ipa-install.html", map[string]string{
		"Name":      filepath.Base(requestPath),
		"PlistLink": plistURL,
	})
}

func parseIPA(path string) (plistInfo *plistBundle, err error) {
	plistRe := regexp.MustCompile(`^Payload/[^/]*/Info\.plist$`)
	reader, err := zip.OpenReader(path)
	defer reader.Close()

	if err != nil {
		return
	}

	var plistFile *zip.File
	for _, file := range reader.File {
		if plistRe.MatchString(file.Name) {
			plistFile = file
			break
		}
	}

	if plistFile == nil {
		err = errors.New("Plist file not found")
		return
	}

	plReader, err := plistFile.Open()
	defer plReader.Close()
	if err != nil {
		return
	}
	buf := make([]byte, plistFile.FileInfo().Size())
	_, err = io.ReadFull(plReader, buf)
	if err != nil {
		return
	}
	dec := goplist.NewDecoder(bytes.NewReader(buf))
	plistInfo = new(plistBundle)
	err = dec.Decode(plistInfo)
	return
}

func generateDownloadPlist(baseURL *url.URL, ipaPath string, plistInfo *plistBundle) ([]byte, error) {
	dp := new(downloadPlist)
	item := new(plItem)

	baseURL.Path = ipaPath
	ipaUrl := baseURL.String()
	item.Assets = append(item.Assets, &plAsset{
		Kind: "software-package",
		URL:  ipaUrl,
	})

	iconFiles := plistInfo.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles
	if iconFiles != nil && len(iconFiles) > 0 {
		baseURL.Path = "/-/unzip/" + ipaPath + "/-/**/" + iconFiles[0] + ".png"
		imgUrl := baseURL.String()
		item.Assets = append(item.Assets, &plAsset{
			Kind: "display-image",
			URL:  imgUrl,
		})
	}

	item.Metadata.Kind = "software"

	item.Metadata.BundleIdentifier = plistInfo.CFBundleIdentifier
	item.Metadata.BundleVersion = plistInfo.CFBundleVersion
	item.Metadata.Title = plistInfo.CFBundleName
	if item.Metadata.Title == "" {
		item.Metadata.Title = filepath.Base(ipaUrl)
	}

	dp.Items = append(dp.Items, item)
	data, err := goplist.MarshalIndent(dp, goplist.XMLFormat, "    ")
	return data, err
}

func combineURL(ctx *gin.Context, path string) *url.URL {
	return &url.URL{
		Scheme: ctx.Request.URL.Scheme,
		Host:   ctx.Request.Host,
		Path:   path,
	}
}

func (s *HTTPStaticServer) genPlistLink(httpPlistLink string) (plistURL string, err error) {
	proxy := s.PlistProxy
	if proxy == "" {
		proxy = "https://plistproxy.herokuapp.com/plist"
	}
	resp, err := http.Get(httpPlistLink)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	retData, err := http.Post(proxy, "text/xml", bytes.NewBuffer(data))
	if err != nil {
		return
	}
	defer retData.Body.Close()

	jsonData, _ := ioutil.ReadAll(retData.Body)
	var ret map[string]string
	if err = json.Unmarshal(jsonData, &ret); err != nil {
		return
	}
	plistURL = proxy + "/" + ret["key"]
	return
}
