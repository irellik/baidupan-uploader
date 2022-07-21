package setting

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Conf struct {
	OAuth struct {
		UrlPrefix    string `yaml:"url_prefix"`
		AppId        string `yaml:"app_id"`
		AppKey       string `yaml:"app_key"`
		SecretKey    string `yaml:"secretKey"`
		RedirectUri  string `yaml:"redirect_uri"`
		AccessToken  string `yaml:"access_token"`
		RefreshToken string `yaml:"refresh_token"`
	}
}

var Cfg *Conf

// 初始化配置
func InitSetting() {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
