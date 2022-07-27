package upload

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/irellik/baidupan-uploader/internal/setting"
	"github.com/irellik/baidupan-uploader/pkg/hashtools"
	openapiclient "github.com/irellik/baidupan-uploader/pkg/openapi"
	"github.com/schollz/progressbar/v3"
)

type PreUploadRequest struct {
	Path      string
	Size      int32
	IsDir     int32
	BlockList []string
	Rtype     int32
}

type PreUploadResponse struct {
	Errno      int64  `json:"errno"`
	RequestId  int64  `json:"request_id"`
	Uploadid   string `json:"uploadid"`
	ReturnType int8   `json:"return_type"`
	BlockList  []int  `json:"block_list"`
}

type PcsUploadRequest struct {
	Path     string
	UploadId string
	Partseq  int
	File     []byte
}

type PcsUploadResponse struct {
	ErrorNo int    `json:"errno"`
	Md5     string `json:"md5"`
}

type FileCreateRequest struct {
	Path      string
	Size      int32
	UploadId  string
	IsDir     int32
	BlockList []string
	Rtype     int32
}

// 预上传
func Precreate(accessToken string, preUploadReq *PreUploadRequest) (preUploadResp *PreUploadResponse, err error) {
	path := url.QueryEscape(preUploadReq.Path)               // string
	isdir := preUploadReq.IsDir                              // int32
	size := preUploadReq.Size                                // int32
	autoinit := int32(1)                                     // int32
	blockListByte, _ := json.Marshal(preUploadReq.BlockList) // string
	blockList := string(blockListByte)
	rtype := preUploadReq.Rtype // int32 | rtype (optional)

	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	_, r, err := api_client.FileuploadApi.Xpanfileprecreate(context.Background()).AccessToken(accessToken).Path(path).Isdir(isdir).Size(size).Autoinit(autoinit).BlockList(blockList).Rtype(rtype).Execute()
	if err != nil {
		return
	}
	// response from `Xpanfileprecreate`: Fileprecreateresponse
	// fmt.Fprintf(os.Stdout, "Response from `FileuploadApi.Xpanfileprecreate`: %v\n", resp)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(bodyBytes, &preUploadResp)
	if preUploadResp.Uploadid == "" {
		err = errors.New(string(bodyBytes))
	}
	return
}

// 分片上传
func MyPcssuperfile2(accessToken string, pcsUploadReq *PcsUploadRequest) (pcsUploadResp *PcsUploadResponse, err error) {
	partseq := strconv.Itoa(pcsUploadReq.Partseq) // string
	path := pcsUploadReq.Path
	uploadid := pcsUploadReq.UploadId // string
	type_ := "tmpfile"
	// 写入文件，为了兼容官方 SDK
	writePath := "/tmp/" + hashtools.HashMd5(pcsUploadReq.File)
	ioutil.WriteFile(writePath, pcsUploadReq.File, 0644)
	defer os.Remove(writePath)
	file, err := os.Open(writePath)
	if err != nil {
		return
	}
	defer file.Close()

	configuration := openapiclient.NewConfiguration()
	//configuration.Debug = true
	api_client := openapiclient.NewAPIClient(configuration)
	_, r, err := api_client.FileuploadApi.Pcssuperfile2(context.Background()).AccessToken(accessToken).Partseq(partseq).Path(path).Uploadid(uploadid).Type_(type_).File(file).Execute()
	if err != nil {
		return
	}
	// response from `Pcssuperfile2`: string
	// fmt.Fprintf(os.Stdout, "Response from `FileuploadApi.Pcssuperfile2`: %v\n", resp)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &pcsUploadResp)

	if pcsUploadResp.Md5 == "" {
		err = errors.New(string(bodyBytes))
	}

	return
}

func FileCreate(accessToken string, fileCreateReq *FileCreateRequest) (resp openapiclient.Filecreateresponse, err error) {
	path := fileCreateReq.Path
	isdir := fileCreateReq.IsDir                              // int32
	size := fileCreateReq.Size                                // int32
	uploadid := fileCreateReq.UploadId                        // string
	blockListByte, _ := json.Marshal(fileCreateReq.BlockList) // string
	blockList := string(blockListByte)
	rtype := fileCreateReq.Rtype // int32 | rtype (optional)

	configuration := openapiclient.NewConfiguration()
	api_client := openapiclient.NewAPIClient(configuration)
	resp, _, err = api_client.FileuploadApi.Xpanfilecreate(context.Background()).AccessToken(accessToken).Path(path).Isdir(isdir).Size(size).Uploadid(uploadid).BlockList(blockList).Rtype(rtype).Execute()
	if err != nil {
		return
	}
	return resp, nil
}

func UploadFile(localPath string, remotePath string) (err error) {
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		log.Fatalln(err)
	}
	blockList := []string{}
	HandleBigFile(localPath, 4*1024*1024, func(index int, shared []byte) {
		blockList = append(blockList, hashtools.HashMd5(shared))
	})
	preCreateResp, err := Precreate(setting.Cfg.OAuth.AccessToken, &PreUploadRequest{
		Path:      remotePath,
		Size:      int32(fileInfo.Size()),
		IsDir:     0,
		BlockList: blockList,
	})
	if err != nil {
		log.Fatalln(err)
	}
	bar := progressbar.Default(int64(len(blockList)))
	HandleBigFile(localPath, 4*1024*1024, func(index int, b []byte) {
		for i := 0; i < 3; i++ {
			_, err := MyPcssuperfile2(setting.Cfg.OAuth.AccessToken, &PcsUploadRequest{
				Path:     remotePath,
				UploadId: preCreateResp.Uploadid,
				Partseq:  index,
				File:     b,
			})
			// log.Println(fmt.Sprintf("uploaded: %d/%d", index+1, len(blockList)))
			bar.Add(1)
			if err == nil {
				break
			}
		}
	})
	for i := 0; i < 3; i++ {
		_, err = FileCreate(setting.Cfg.OAuth.AccessToken, &FileCreateRequest{
			Path:      remotePath,
			Size:      int32(fileInfo.Size()),
			IsDir:     0,
			BlockList: blockList,
			UploadId:  preCreateResp.Uploadid,
			Rtype:     3,
		})
		if err == nil {
			break
		}
	}

	return
}
