package common

import (
	"context"
	"errors"
	"fmt"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/dialer"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/log"
	taskPb "git-biz.360es.cn/zion-infra/file_access-api/file_access/task"
	apiPb "git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api"
	//logApiPb "git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/log"
	storageApiPb "git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/storage"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MAX_ERROR_TIMES = 2

type UploadSrv struct {
	storgeCli apiPb.StorageClient
	timeout   time.Duration
}

func NewUploadSrv() *UploadSrv {
	timeout := 2 * time.Minute
	fileAccessConn, err := dialer.Dial(context.Background(), SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[UploadSrv]connect py rpc error:%v", err)
		return nil

	}
	storageCli := apiPb.NewStorageClient(fileAccessConn)
	return &UploadSrv{storgeCli: storageCli, timeout: timeout}
}

func (s *UploadSrv) UploadFile(filePath string) (string, error) {
	// 首先读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", errors.New("read file content error,err:" + err.Error())
	}
	sha1 := GetSha1(data)
	return s.UploadFileContent(fmt.Sprintf("files/%s", sha1), data)
}

func (s *UploadSrv) UploadFileContent(name string, fileContent []byte) (string, error) {
	//step1 : 初始化上传file
	log.Infof("upload file content begin")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	req := CreateStorageRequest()
	req.Parent = name
	req.Storage.Uri = fmt.Sprintf("%s/storage", name)
	req.Storage.File.Size = int64(len(fileContent))
	resp, err := s.storgeCli.CreateStorage(timeoutCtx, req)
	if err != nil {
		return "", errors.New("request file create storage err:" + err.Error())
	}
	//调用create storage 得到一个url
	logUri := resp.Uri

	//step2 : 获取上传url任务
	var resultErr error
	for retryTimes := 0; retryTimes < MAX_ERROR_TIMES; retryTimes++ {
		tasks, err := s.initFetchTask(logUri)
		if err != nil {
			return "", errors.New("request init fetch tasks error")
		} else {
			if len(tasks) == 0 {
				return name, nil
			}
		}
		resultErr := s.runTask(tasks, logUri, fileContent)
		if resultErr != nil {
			continue
		}
		break
	}

	return name, resultErr
}

func (s *UploadSrv) initFetchTask(uri string) ([]*taskPb.FileUploadTask, error) {
	log.Infof("fetch task start")
	tasks := make([]*taskPb.FileUploadTask, 0)
	taskReq := &storageApiPb.ListUploadTasks_Request{Uri: uri}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	// 根据uri拉取 tasks
	taskResp, err := s.storgeCli.ListUploadTasks(timeoutCtx, taskReq)
	if err != nil {
		return tasks, errors.New("get upload file url error")
	}
	uploadTasks := taskResp
	// 得到待下载的tasks列表
	tasks = uploadTasks.GetWaitingTasks()
	log.Infof("fetch task end")
	return tasks, nil
}

func (s *UploadSrv) runTask(tasks []*taskPb.FileUploadTask, name string, fileContents []byte) error {

	log.Infof("run task start,tasks lens %d,start_time", strconv.Itoa(len(tasks)), time.Now().Unix())
	var wg sync.WaitGroup
	var errData error = nil
	for i, task := range tasks {
		wg.Add(1)
		go func(task *taskPb.FileUploadTask, i int) {
			defer wg.Done()
			log.Infof("upload file start,time", time.Now().Unix())
			err := s.exeTask(task.Url, fileContents[task.Offset:task.Length+task.Offset])
			log.Infof("finish upload file end,time", time.Now().Unix())
			if err != nil {
				task.State = taskPb.State_FAILED
			} else {
				task.State = taskPb.State_FINISHED
			}

			taskReq := &storageApiPb.ListUploadTasks_Request{Uri: name}
			taskReq.FinishedTasks = []*taskPb.FileUploadTask{task}
			timeoutCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
			defer cancel()
			taskResp, err := s.storgeCli.ListUploadTasks(timeoutCtx, taskReq)
			log.Infof("get new upload task end,time", time.Now().Unix())
			if err != nil {
				errData = err
			} else {
				newTasks := taskResp.GetWaitingTasks()
				if len(newTasks) > 0 {
					err := s.runTask(newTasks, name, fileContents)
					if err != nil {
						errData = err
					}
				}
			}
		}(task, i)
	}
	wg.Wait()
	log.Infof("run task end")
	return errData
}

func (s *UploadSrv) exeTask(url string, datas []byte) error {
	request, err := http.NewRequest("PUT", url, strings.NewReader(string(datas)))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 && res.StatusCode != 100 {
		return errors.New("http response code is error,code:")
	}

	return nil
}

func (s *UploadSrv) exeDownload(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("http response code is error,err:" + err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("read all datas error,err:" + err.Error())
	}
	return body, nil
}

func (s *UploadSrv) DownloadFileUrl(name string) (string, error) {
	req := &storageApiPb.GetStorage_Request{Uri: name}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	resp, err := s.storgeCli.GetStorage(timeoutCtx, req)
	if err != nil {
		return "", errors.New("get file download url error,err:" + err.Error())
	}
	if resp.File.Url == "" {
		return "", errors.New("get file download url error,download url is empty")
	}

	return resp.File.Url, nil
}

func (s *UploadSrv) ReadFileDatas(name string) ([]byte, error) {
	url, err := s.DownloadFileUrl(name)
	if err != nil {
		return nil, errors.New("get file download url error,err:" + err.Error())
	}

	datas, err := s.exeDownload(url)
	if err != nil {
		return nil, errors.New("get file data error")
	}

	return datas, nil
}

func (s *UploadSrv) DownloadFile(name, downloadPath string) (string, error) {
	url, err := s.DownloadFileUrl(name)
	if err != nil {
		return "", errors.New("get file download url error,err:" + err.Error())
	}

	datas, err := s.exeDownload(url)
	if err != nil {
		return url, errors.New("get file data error,err:" + err.Error())
	}

	var f *os.File
	if checkFileIsExist(downloadPath) { //如果文件存在
		f, err = os.OpenFile(downloadPath, os.O_WRONLY, 0666) //打开文件
	} else {
		f, err = os.Create(downloadPath) //创建文件
	}
	if err != nil {
		return url, errors.New("open file error,err:" + err.Error())
	}

	_, err = io.WriteString(f, string(datas)) //写入文件(字符串)
	if err != nil {
		return url, errors.New("save file data error,err:" + err.Error())
	}
	return url, nil
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
