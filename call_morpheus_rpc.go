package common

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/dialer"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/log"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/file_attribute"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/file_detection"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/storage"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/cert"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/file"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/protobuf/field_mask"
)

// 此文件用于morpheus rpc客户端封装

// 从morpheus上得到s3的下载地址
func GetS3UrlByMorpheusRpc(ctx context.Context, originUrl string, key string, ch chan<- map[string]interface{}) error {
	log.Debugf("Entring GetS3UrlByMorpheusRpc")
	storageReq := new(storage.GetStorage_Request)
	storageReq.Uri = fmt.Sprintf("%s/storage", originUrl)
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS_STORAGE)
	keylToS3Url := make(map[string]interface{}, 0)
	if err != nil {
		log.Errorf("[GetS3UrlByMorpheusRpc]connect py rpc error:%v", err)
		//keylToS3Url[key] = fmt.Sprintf("http://127.0.0.1/%s.json", key)
		keylToS3Url[key] = nil
		ch <- keylToS3Url
		return err

	}
	morpheusApi := api.NewStorageClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()
	storageRes, err := morpheusApi.GetStorage(ctx, storageReq)
	log.Debugf("GetS3UrlByMorpheusRpc request :%+v  res: %+v", storageReq, storageRes)
	var s3Url interface{}
	s3Url = nil
	if err != nil {
		log.Infof("GetStorage originUrl:%v get s3_url error:%v", originUrl, err)
	} else {
		s3Url = storageRes.File.Url
	}

	// chan中的每一个元素 存储着key到s3url的单个映射,

	keylToS3Url[key] = s3Url
	ch <- keylToS3Url

	return nil

}

// TODO 此函数暂时存在，不请求 过完typer qex 拿到结果后，是update还是create morpheus
func FileAttributeCreateFileByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, typer *TyperInfo, qex *QexInfo) (*file.File, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[FileAttributeUpdateFileByMorpheus]connect py rpc error:%v", err)
		return nil, err

	}
	morpheusApi := api.NewFileAttributeClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()
	//request := new(file_attribute.CreateFile_Request)
	request := CreateFileAttributeRequest()
	sha1 := fileBaseInfo.Sha1
	size := fileBaseInfo.Size
	hash := typer.Hash
	md5 := hash.Md5
	sha256 := hash.Sha256
	ssdeep := hash.Ssdeep
	request.File.Uri = fmt.Sprintf("%s/%s", "files", sha1)
	request.File.Size = uint64(size)
	request.File.Md5 = md5
	request.File.Sha256 = sha256
	request.File.State = file.File_NEW_CREATED
	request.File.Ssdeep = ssdeep
	// request.File.Type = file.FileType(0)
	// 格式:pe_exe_x86
	file_type_by_qex := qex.FileType
	trid := typer.Trid
	exiftool := typer.Exiftool
	judgeType := JudgeFileType(file_type_by_qex, trid, exiftool)
	// 转化成pb中定义的枚举文件类型
	fileType, ok := file.FileType_value[judgeType]
	if !ok {
		fileType = 0
	}
	request.File.Type = file.FileType(fileType)

	// TODO 这两个未定
	request.File.Vhash = "NOIMPLEMENT"
	request.File.Hash_360 = "NOIMPLEMENT"

	//request.Mask.Paths = []string{"file.uri", "file.size", "file.md5", "file.sha256", "file.state", "file.ssdeep", "file.type"}
	res, err := morpheusApi.CreateFile(ctx, request)
	if err != nil {
		log.Warnf("FileAttributeCreateFileByMorpheus error:%v", err)
		return nil, err
	}
	return res, nil

}

func FileAttributeUpdateFileByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, typer *TyperInfo, qex *QexInfo) (*file.File, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[FileAttributeUpdateFileByMorpheus]connect py rpc error:%v", err)
		return nil, err

	}
	morpheusApi := api.NewFileAttributeClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()
	//request := new(file_attribute.CreateFile_Request)
	request := UpdateFileAttributeRequest()
	sha1 := fileBaseInfo.Sha1
	size := fileBaseInfo.Size
	hash := typer.Hash
	md5 := hash.Md5
	sha256 := hash.Sha256
	ssdeep := hash.Ssdeep
	request.File.Uri = fmt.Sprintf("%s/%s", "files", sha1)
	request.File.Size = uint64(size)
	request.File.Md5 = md5
	request.File.Sha256 = sha256
	//request.File.State = file.File_UPLOAD_COMPLETE
	request.File.Ssdeep = ssdeep
	// request.File.Type = file.FileType(0)
	// 格式:pe_exe_x86
	file_type_by_qex := qex.FileType
	trid := typer.Trid
	exiftool := typer.Exiftool
	judgeType := JudgeFileType(file_type_by_qex, trid, exiftool)
	// 转化成pb中定义的枚举文件类型
	fileType, ok := file.FileType_value[judgeType]
	if !ok {
		fileType = 0
	}
	request.File.Type = file.FileType(fileType)

	// TODO 这两个未定
	//request.File.Vhash = "NOIMPLEMENT"
	//request.File.Hash_360 = "NOIMPLEMENT"

	request.Mask.Paths = []string{"Uri", "Size", "Md5 ", "Sha256", "Ssdeep", "Type"}
	res, err := morpheusApi.UpdateFile(ctx, request)
	log.Debugf("FileAttributeUpdateFileByMorpheus Step 5-1 resquest:%+v res: %+v\n", request, res)
	if err != nil {
		log.Warnf("FileAttributeUpdateFileByMorpheus error:%v", err)
		return nil, err
	}
	return res, nil

}

//初次创建static属性
func FileAttributeCreateStaticByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, typer *TyperInfo, qex *QexInfo) (*file.Static, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[FileAttributeUpdateFileByMorpheus]connect py rpc error:%v", err)
		return nil, err

	}
	morpheusApi := api.NewFileAttributeClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()
	sha1 := fileBaseInfo.Sha1
	exiftool := typer.Exiftool
	icon := typer.Icon
	staticReq := new(file_attribute.GetStatic_Request)
	staticReq.Uri = fmt.Sprintf("%s/%s/static", "files", sha1)
	static, err := morpheusApi.GetStatic(ctx, staticReq)
	if err != nil {
		log.Warnf("GetStatic Step 6-3 error %+v", err)

	}
	log.Debugf("Getstatic Step 6-4 request:%+v res:%+v", staticReq, static)
	file_type_by_qex := qex.FileType
	trid := typer.Trid
	judgeType := strings.ToLower(JudgeFileType(file_type_by_qex, trid, exiftool))
	desc, ok := EXT_DESC[judgeType]
	typer_cert := typer.Cert
	certs := make([]*cert.Certificate, 0)

	// 处理证书逻辑
	for _, v := range typer_cert {

		cert := CreateCertificateRequest()
		cert.Sha1 = v.Sha1
		cert.Name = v.Organization

		t, _ := time.Parse(TIME_FORMAT, v.StartTime)

		cert.ValidFrom.Seconds = t.Unix()
		t, _ = time.Parse(TIME_FORMAT, v.EndTime)
		cert.ValidTo.Seconds = t.Unix()

		cert.SerialNumber = v.SerialNumber
		certs = append(certs, cert)
	}
	exifMap := make(map[string]string)
	for _, v := range exiftool {
		_v := strings.Split(v, ":")
		exifMap[strings.Trim(_v[0], " ")] = strings.Trim(_v[1], " ")
	}

	if err != nil {
		// static不存在
		request := CreateStaticRequest()
		request.Parent = fmt.Sprintf("files/%s", sha1)
		request.Static.Uri = fmt.Sprintf("files/%s/static", sha1)
		if ok {
			request.Static.Magic = desc
		}
		request.Static.StaticFormats = []string{"NOIMPLEMENT"}
		request.Static.Executable.Pe.Signature.Certs = certs
		if icon.Png != "" {
			request.Static.Executable.Pe.MainIcon.RawBase64 = icon.Png
		}
		request.Static.ExiftoolMeta = exifMap
		res, err := morpheusApi.CreateStatic(ctx, request)
		log.Debugf("FileAttributeCreateStaticByMorpheus Step 6-1 request %+v  res:%+v\n", request, res)
		if err != nil {
			log.Warnf("FileAttributeCreateStaticByMorpheus error:%v", err)
			return nil, err
		}
		return res, nil

	} else {
		// static 存在
		if ok {
			static.Magic = desc
		}
		static.StaticFormats = []string{"NOIMPLEMENT"}
		if static.Executable.Pe.GetSignature() == nil {
			static.Executable.Pe.Signature = new(file.Executable_Signature)
		}
		if static.Executable.Pe.GetMainIcon() == nil {
			static.Executable.Pe.MainIcon = new(file.Executable_PE_MainIcon)
		}
		static.Executable.Pe.Signature.Certs = certs
		if icon.Png != "" {
			static.Executable.Pe.MainIcon.RawBase64 = icon.Png
		}
		static.ExiftoolMeta = exifMap
		request := new(file_attribute.UpdateStatic_Request)
		mark := new(field_mask.FieldMask)
		mark.Paths = []string{"Executable"}
		request.Static = static
		request.Mask = mark
		res, err := morpheusApi.UpdateStatic(ctx, request)
		if err != nil {
			log.Errorf("Step 6-6 FileAttributeUpdateStaticByMorpheus error:%v", err)
			return nil, err
		}
		log.Debugf("Step 6-7 Want to create  but update static request:%+v res:%+v", request, res)
		return res, nil

	}

}

func FileAttributeUpdateStaticByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, owl *OwlInfo) (*file.Static, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[FileAttributeUpdateStaticByMorpheus]connect py rpc error:%v", err)
		return nil, err
	}
	morpheusApi := api.NewFileAttributeClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()

	sha1 := fileBaseInfo.Sha1
	uri := fmt.Sprintf("%s/%s/static", "files", sha1)
	staticReq := new(file_attribute.GetStatic_Request)
	staticReq.Uri = uri
	static, err := morpheusApi.GetStatic(ctx, staticReq)
	log.Debugf("@@@Before %+v", static)
	if err != nil {
		log.Warnf("FileAttributeGetStaticByMorpheus uir:%v, error:%v", uri, err)
		return nil, err
	}
	stream := owl.Streams[0]
	request := new(file_attribute.UpdateStatic_Request)
	log.Debugf("FILETYPE Step 9-1 %v", stream.FType)
	if stream.FType == "pe" {
		isX64 := stream.Bit
		var executable *file.Executable
		if static.Executable == nil {
			executable = new(file.Executable)
			executable.Uri = uri
		} else {
			executable = static.Executable
		}
		peExecutable := new(file.Executable_PE)
		if isX64 == 0 {
			peExecutable.IsX64 = false
		} else {
			peExecutable.IsX64 = true
		}
		peExecutable.PdbPath = stream.Pdb
		peExecutable.CompilerType = stream.CType
		version := stream.Versions[0]
		fileVersion := new(file.Executable_PE_FileVersion)
		fileVersion.CompanyName = version.CN
		fileVersion.LegalCopyright = version.LC
		fileVersion.InternalName = version.IN
		fileVersion.FileVersion = version.FV
		fileVersion.ProductVersion = version.PV
		fileVersion.OriginalFileName = version.ON
		peExecutable.FileVersion = fileVersion
		header := new(file.Executable_PE_Header)
		ct := time.Unix(stream.Header.Time, 0)
		header.CompilationTime = ct.Format(TIME_FORMAT)
		peExecutable.Header = header
		sts := make([]*file.Executable_PE_Section, 0)
		for _, item := range stream.Header.Sections {
			section := new(file.Executable_PE_Section)
			section.Name = item.Name
			sts = append(sts, section)
		}
		peExecutable.Sections = sts
		crs := make([]*file.Executable_PE_ContainedResource, 0)
		for _, item := range stream.Resources {
			res := new(file.Executable_PE_ContainedResource)
			res.Type = item.SType
			crs = append(crs, res)
		}
		peExecutable.Resources = crs
		sepi := make([]*file.Executable_PE_Import, 0)
		for _, item := range stream.Imports {
			its := new(file.Executable_PE_Import)
			its.ModuleName = item.Name
			its.FunctionName = item.Apis
			sepi = append(sepi, its)
		}
		peExecutable.Imports = sepi
		sepe := make([]*file.Executable_PE_Export, 0)
		ite := new(file.Executable_PE_Export)
		ite.ModuleName = stream.Export.Name
		ite.FunctionName = stream.Export.Apis
		sepe = append(sepe, ite)
		peExecutable.Exports = sepe
		executable.Pe = peExecutable
		static.Executable = executable
		request.Static = static
		log.Debugf("@@@After %+v", request)
	} else {
		log.Debug("not pe file return ")
		return nil, nil

	}
	mark := new(field_mask.FieldMask)
	mark.Paths = []string{"Executable"}
	request.Mask = mark

	res, err := morpheusApi.UpdateStatic(ctx, request)
	log.Debugf("FileAttributeUpdateStaticByMorpheus Step 9-2 request: %+v  res:%+v \n", request, res)
	if err != nil {
		log.Errorf("FileAttributeUpdateStaticByMorpheus error:%v", err)
		return nil, err
	}
	return res, nil
}

func FileDetectionCreateByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, fileDetections map[string]interface{}) (*file.Detection, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS_DETECTION)
	if err != nil {
		log.Errorf("[FileDetectionFileByMorpheus]connect py rpc error:%v", err)
		return nil, err
	}
	morpheusApi := api.NewFileDetectionClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()

	sha1 := fileBaseInfo.Sha1
	timeStamp := int64(time.Now().Second())
	request := new(file_detection.CreateDetection_Request)
	request.Parent = fmt.Sprintf("%s/%s", "files", sha1)
	detection := new(file.Detection)
	detection.Uri = fmt.Sprintf("%s/%s/detections", "files", sha1)
	detection.Sha1 = sha1
	results := make([]*file.Detection_Result, 0)
	for avName, avInfo := range fileDetections {
		avResult := new(file.Detection_Result)
		avResult.Engine = avName
		avResult.DetectTime = new(timestamp.Timestamp)
		avResult.DetectTime.Seconds = timeStamp
		avResult.Uri = fmt.Sprintf("files/%s/detections/%s", sha1, avName)
		switch avName {
		case "qex":
			qexInfo, _ := avInfo.(*QexInfo)
			avResult.MalwareName = qexInfo.MalwareName
			break
		case "bole":
			qvmInfo, _ := avInfo.(*QvmInfo)
			avResult.MalwareName = qvmInfo.MalwareName
			break
		case "ave":
			aveInfo, _ := avInfo.(*AveInfo)
			avResult.MalwareName = aveInfo.MalwareName
			break
		case "bd":
			bdInfo, _ := avInfo.(*BdInfo)
			avResult.MalwareName = bdInfo.MalwareName
			break
		case "sign":
			signInfo, _ := avInfo.(*SignInfo)
			if signInfo.Status == 1 {
				avResult.MalwareName = signInfo.MalwareName
			}
			break
		case "owl":
			owlInfo, _ := avInfo.(*OwlInfo)
			if len(owlInfo.OwlTags) > 0 {
				avResult.MalwareName = owlInfo.OwlTags[0]
			}
			break
		}
		results = append(results, avResult)
	}
	detection.Results = results
	request.Detection = detection

	res, err := morpheusApi.CreateDetection(ctx, request)
	log.Debugf("FileDetectionCreateByMorpheus Step 7/10-1 request :%+v res: %+v\n", request, res)
	if err != nil {
		log.Warnf("FileDetectionCreateByMorpheus error:%v", err)
		return nil, err
	}
	return res, nil
}
func FileBehaviorCreateByMorpheus(ctx context.Context, fileBaseInfo FileBaseInfo, logPath string) (*file.Behavior, error) {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[FileDetectionFileByMorpheus]connect py rpc error:%v", err)
		return nil, err
	}
	sha1 := fileBaseInfo.Sha1
	logPath_ := filepath.Dir(logPath)
	pcapSha1, reportJsonSha1, snapSha1s := GetPcapAndSnapShot(logPath_)
	pcapSha1 = fmt.Sprintf("files/%s/behaviors/cuckoo/%s", sha1, pcapSha1)

	snapShotSha1s := make([]string, 0)
	for _, k := range snapSha1s {
		snapShotSha1s = append(snapShotSha1s, fmt.Sprintf("files/%s/behaviors/cuckoo/%s", sha1, k))
	}

	req := CreateBehaviorRequest()

	req.Parent = fmt.Sprintf("files/%s", sha1)
	req.Behavior.Uri = fmt.Sprintf("files/%s/behaviors", sha1)
	req.Behavior.CuckooBehavior.Uri = fmt.Sprintf("files/%s/behaviors/cuckoo", sha1)
	req.Behavior.CuckooBehavior.Sha1 = sha1
	req.Behavior.CuckooBehavior.PcapUri = pcapSha1
	req.Behavior.CuckooBehavior.SnapshotUris = snapShotSha1s
	req.Behavior.CuckooBehavior.ReportUri = fmt.Sprintf("files/%s/behaviors/cuckoo/%s", sha1, reportJsonSha1)
	morpheusApi := api.NewFileBehaviorClient(conn)
	res, err := morpheusApi.CreateBehavior(ctx, req)
	log.Debugf("[FileBehaviorCreateByMorpheus] CreateBehavior request :%+v  res: %+v ", req, res)
	if err != nil {
		log.Warnf("FileBehaviorCreateByMorpheus error:%v", err)
		return nil, err
	}
	return res, nil

}

//更新file的状态
func UpdateFileStateByMorpheus(ctx context.Context, sha1, state string) error {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_MORPHEUS)
	if err != nil {
		log.Errorf("[UpdateFileStateByMorpheus]connect morpheus rpc error:%v", err)
		return err

	}
	morpheusApi := api.NewFileAttributeClient(conn)
	request := UpdateFileAttributeRequest()
	request.File.Uri = fmt.Sprintf("files/%s", sha1)
	FileState, ok := file.File_State_value[state]
	if ok == false {
		request.File.State = file.File_UNSPECIFIED
	}
	request.File.State = file.File_State(FileState)
	request.Mask.Paths = []string{"State"}
	res, err := morpheusApi.UpdateFile(ctx, request)
	log.Debugf("[UpdateFileStateByMorpheus] request %+v ,res %+v", request, res)
	if err != nil {
		log.Warnf("UpdateFileStateByMorpheus error: %v", err)
		return err
	}
	return nil
}
