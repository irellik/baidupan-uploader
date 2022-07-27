# baidupan-uploader

不用暴露百度账号上传文件的小工具，使用百度网盘开发者 + oAuth，目前需要自己申请成为百度网盘开发者，并创建一个应用。

## 准备

1. 在[百度网盘开发平台](https://pan.baidu.com/union/home)认证成为开发者
2. 在控制台创建一个应用，获取到应用的 AppId、AppKey、SecretKey、SignKey
3. 将信息填写到运行目录 config.yaml 中

## 使用
### 常规
假设编译后的二进制文件名为 upload，执行 upload /local/path /baidupan/path
### docker
docker run -it -v /config/path:/opt/app -v /local/path:/local/path iwww/baiduupload:1.1.0 /local/path /baidupan/path