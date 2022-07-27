package oauth

import (
	"log"
	"net/url"

	"github.com/irellik/baidupan-uploader/internal/setting"
)

func GetOauthUrl() string {
	base, err := url.Parse(setting.Cfg.OAuth.UrlPrefix)
	if err != nil {
		log.Fatal(err)
	}

	base.Path += "authorize"

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", setting.Cfg.OAuth.AppKey)
	params.Add("redirect_uri", setting.Cfg.OAuth.RedirectUri)
	params.Add("scope", "basic,netdisk")
	params.Add("device_id", setting.Cfg.OAuth.AppId)

	base.RawQuery = params.Encode()
	return base.String()
}
