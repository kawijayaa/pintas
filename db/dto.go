package db

type ShortUrlDto struct {
	Path        string `form:"path"`
	RedirectUrl string `form:"redirect_url"`
}
