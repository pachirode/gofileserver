package v1

type GetUserResponse UserInfo

type UserInfo struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type CreateUserRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" valid:"required,stringlength(1|255)"`
	Email    string `json:"email" valid:"required,email"`
}

type ListUserRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ListUserResponse struct {
	TotalCount int64       `json:"totalCount"`
	Users      []*UserInfo `json:"users"`
}
