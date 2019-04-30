package common

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/log"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// 工具类的辅助函数
func IsInArray(item interface{}, array interface{}) bool {
	targetValue := reflect.ValueOf(array)
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == item {
				return true
			}
		}
	}
	return false
}

// url 是 s3的下载地址，name 是workload中的name
func DownLogByUrl(url, name string) (string, error) {
	log.Debugf("DownLogByUrl begin to download %v", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Warnf("http get url:%v error:%v", url, err)
		return "", err
	}
	defer resp.Body.Close()
	_url := strings.Split(url, "/")

	logPath := filepath.Join(DOWN_LOAD_PATH, name)
	fileName := filepath.Join(logPath, GetSha1([]byte(_url[len(_url)-1:][0])))
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		oldMask := syscall.Umask(0)
		err = os.MkdirAll(logPath, os.ModePerm)
		syscall.Umask(oldMask)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Warnf("read body error:%v", err)
		return "", err
	}
	err = ioutil.WriteFile(fileName, body, 0666)
	log.Debugf("DownLogByUrl origin_url:%v local_path %+v", url, fileName)
	return fileName, nil
}

// 返回文件类型
// " 82.0% (.EXE) Win64 Executable (generic) (27624/17/4)"
//
func JudgeFileType(file_type_by_qex string, trid, exiftool []string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Infof("catch error:%v", PanicTrace(1))

		}
	}()
	filerArray := []string{"c", "html", "unicode", "bin", "mime"}
	fileType := strings.Split(file_type_by_qex, "_")
	_type := fileType[0]
	length := len(fileType)
	if length > 1 {
		_type = fileType[1]
	}
	flag := IsInArray(_type, filerArray)
	// 不在filterArray中的类型，直接就可依据qex中的file_type字段确定文件的类型
	if flag == false {
		if length > 1 {
			return strings.ToUpper(strings.Join(fileType[0:2], ""))
		} else {
			return strings.ToUpper(fileType[0])

		}

	} else {
		tridFirst := trid[0]
		//a="82.0"
		percent, _ := strconv.ParseFloat(strings.Trim(strings.Split(tridFirst, "%")[0], " "), 32)
		start := strings.IndexAny(tridFirst, "(")
		end := strings.IndexAny(tridFirst, ")")
		fileType := tridFirst[start+2 : end]
		if percent >= 60 {
			_res := strings.ToUpper(fileType)
			if IsInArray(_res, []string{"EXE", "DLL"}) == true {
				return strings.ToUpper("PE" + _res)

			}
			return _res

		} else {
			exifMap := make(map[string]string)
			for _, v := range exiftool {
				_v := strings.Split(v, ":")
				exifMap[strings.Trim(_v[0], " ")] = strings.Trim(_v[1], " ")
			}
			Extension, _ := exifMap["File Type Extension"]
			if strings.ToUpper(strings.Trim(fileType, " ")) == strings.ToUpper(strings.Trim(Extension, " ")) {
				_res := strings.ToUpper(fileType)
				if IsInArray(_res, []string{"EXE", "DLL"}) == true {
					return strings.ToUpper("PE" + _res)

				}
				return _res
			} else {
				return "UNSPECIFIED"
			}

		}
	}

}

func GetSha1(f []byte) string {
	h := sha1.New()
	h.Write(f)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func GetSha1ByPath(p string) string {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		log.Warnf("GetSha1ByPath path: %v error:%v", p, err)
		return ""
	}
	return GetSha1(content)
}

func UnzipFile(filePath string) string {
	dir := filepath.Dir(filePath)
	cmd := exec.Command("unzip", "-d", dir, filePath)
	cmd.Run()
	return filepath.Join(dir, "latest")

}

func GetUid() string {
	uid := uuid.NewV4()
	suid := strings.Replace(fmt.Sprintf("%s", uid), "-", "", -1)
	return suid
}

func PanicTrace(kb int) string {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, kb<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	return string(stack)
}
