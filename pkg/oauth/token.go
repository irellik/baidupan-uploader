package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/irellik/baidupan-uploader/internal/setting"
	"gopkg.in/yaml.v3"

	openapiclient "github.com/irellik/baidupan-uploader/pkg/openapi"
)

type TokenResponse struct {
	ExpiresIn     int32  `json:"expires_in"`
	RefreshToken  string `json:"refresh_token"`
	AccessToken   string `json:"access_token"`
	SessionSecret string `json:"session_secret,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	Scope         string `json:"scope"`
}

type UserInfo struct {
	AvatarUrl string `json:"avatar_url"`
	Name      string `json:"baidu_name"`
	ErrorMsg  string `json:"errmsg"`
	ErrorNo   int64  `json:"errno"`
	UK        int64  `json:"uk"`
	VipType   int64  `json:"vip_type"`
}

// 通过 code 获取 token
func getTokenFromCode(code string) (tokenResp TokenResponse, err error) {
	// 构造请求参数
	req, err := http.NewRequest("GET", setting.Cfg.OAuth.UrlPrefix, nil)
	if err != nil {
		return
	}
	req.URL.Path += "token"
	query := req.URL.Query()
	query.Add("grant_type", "authorization_code")
	query.Add("code", code)
	query.Add("redirect_uri", setting.Cfg.OAuth.RedirectUri)
	query.Add("client_id", setting.Cfg.OAuth.AppKey)
	query.Add("client_secret", setting.Cfg.OAuth.SecretKey)
	req.URL.RawQuery = query.Encode()

	// 发起请求
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &tokenResp)
	// 如果非预期返回 err
	if tokenResp.ExpiresIn == 0 {
		err = errors.New(string(body))
	}
	return
}

// 刷新 access token
func refreshToken(refresh_token string) (tokenResp TokenResponse, err error) {
	// 构造请求参数
	req, err := http.NewRequest("GET", setting.Cfg.OAuth.UrlPrefix, nil)
	if err != nil {
		return
	}
	req.URL.Path += "token"
	query := req.URL.Query()
	query.Add("grant_type", "refresh_token")
	query.Add("refresh_token", refresh_token)
	query.Add("client_id", setting.Cfg.OAuth.AppKey)
	query.Add("client_secret", setting.Cfg.OAuth.SecretKey)
	req.URL.RawQuery = query.Encode()

	// 发起请求
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 解析响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &tokenResp)
	// 如果非预期返回 err
	if tokenResp.ExpiresIn == 0 {
		err = errors.New(string(body))
	}
	return
}

// 写入 token
func writeToken(accessToken string, refreshToken string) (err error) {
	setting.Cfg.OAuth.AccessToken = accessToken
	setting.Cfg.OAuth.RefreshToken = refreshToken
	byteList, err := yaml.Marshal(setting.Cfg)
	if err != nil {
		return
	}
	err = ioutil.WriteFile("config.yaml", byteList, 0644)
	if err != nil {
		return
	}
	return nil
}

// 获取用户数据
func getUserInfo(accessToken string) (user UserInfo, err error) {
	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	_, r, err := api_client.UserinfoApi.Xpannasuinfo(context.Background()).AccessToken(accessToken).Execute()
	if err != nil {
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &user)
	if user.ErrorMsg != "succ" {
		err = errors.New(string(bodyBytes))
	}
	return
}

// 初始化
func InitToken() (err error) {
	// 尝试使用 access token 获取信息
	_, err = getUserInfo(setting.Cfg.OAuth.AccessToken)
	if err == nil {
		return
	}
	// 尝试使用 refresh token 获取 access token
	tokenResp, err := refreshToken(setting.Cfg.OAuth.RefreshToken)
	if err == nil {
		writeToken(tokenResp.AccessToken, tokenResp.RefreshToken)
		return
	}
	// 重新获取 access token
	var code string
	for {
		fmt.Println(fmt.Sprintf("请访问该链接获取到 code 后，输入该 code 后回车：%s", GetOauthUrl()))
		fmt.Scanln(&code)
		tokenResp, err = getTokenFromCode(code)
		if err == nil {
			writeToken(tokenResp.AccessToken, tokenResp.RefreshToken)
			break
		}
	}
	return
}
