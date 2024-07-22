package db

type UrlDto struct {
	Path        string `form:"path" json:"path"`
	RedirectUrl string `form:"redirect_url" json:"redirect_url"`
}

type UserDto struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}
