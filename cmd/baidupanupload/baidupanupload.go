package main

import (
	"log"
	"os"

	"github.com/irellik/baidupan-uploader/internal/setting"
	"github.com/irellik/baidupan-uploader/pkg/oauth"
	"github.com/irellik/baidupan-uploader/pkg/upload"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	setting.InitSetting()
	oauth.InitToken()
	if len(os.Args) < 2 {
		log.Fatalln("请指定要上传的文件")
	}
	if len(os.Args) < 3 {
		log.Fatalln("请指定要上传的目标路径")
	}
	// 获取文件大小
	localPath := os.Args[1]
	remotePath := os.Args[2]
	err := upload.UploadFile(localPath, remotePath)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println("file uploaded!")
}
