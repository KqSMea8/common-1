package common

import (
	ffile "git-biz.360es.cn/zion-infra/file_access-api/file_access/file"
	k_api "git-biz.360es.cn/zion-infra/kamala-api/kamala/api/job"
	kamala_job "git-biz.360es.cn/zion-infra/kamala-api/kamala/job"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/file_attribute"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/file_behavior"
	storageApiPb "git-biz.360es.cn/zion-infra/morpheus-api/morpheus/api/storage"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/cert"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/file"
	"git-biz.360es.cn/zion-infra/morpheus-api/morpheus/storage"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/protobuf/api"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
)

// 使用时直接调用，以免嵌套结构到处new
func CreateKamakaJobRequest() *k_api.CreateJob_Request {
	jobRequest := new(k_api.CreateJob_Request)
	jobRequest.Job = new(kamala_job.Job)
	jobRequest.Job.CompleteJobCallback = new(api.Api)
	return jobRequest
}

func CreateFileAttributeRequest() *file_attribute.CreateFile_Request {
	request := new(file_attribute.CreateFile_Request)
	request.File = new(file.File)
	request.File.Community = new(file.Community)
	request.File.Submission = new(file.Submission)
	return request

}

func UpdateFileAttributeRequest() *file_attribute.UpdateFile_Request {
	request := new(file_attribute.UpdateFile_Request)
	request.File = new(file.File)
	request.File.Community = new(file.Community)
	request.File.Submission = new(file.Submission)
	request.Mask = new(field_mask.FieldMask)
	return request

}

func CreateStaticRequest() *file_attribute.CreateStatic_Request {
	request := new(file_attribute.CreateStatic_Request)
	request.Static = new(file.Static)
	request.Static.Executable = new(file.Executable)
	request.Static.Executable.Pe = new(file.Executable_PE)
	request.Static.Executable.Pe.Signature = new(file.Executable_Signature)
	request.Static.Executable.Pe.MainIcon = new(file.Executable_PE_MainIcon)
	return request
}

func CreateBehaviorRequest() *file_behavior.CreateBehavior_Request {
	request := new(file_behavior.CreateBehavior_Request)
	request.Behavior = new(file.Behavior)
	request.Behavior.CuckooBehavior = new(file.CuckooBehavior)
	return request
}

func CreateCertificateRequest() *cert.Certificate {
	request := new(cert.Certificate)
	request.ValidFrom = new(timestamp.Timestamp)
	request.ValidTo = new(timestamp.Timestamp)
	return request
}

func CreateStorageRequest() *storageApiPb.CreateStorage_Request {
	request := new(storageApiPb.CreateStorage_Request)
	request.Storage = new(storage.Storage)
	request.Storage.File = new(ffile.File)
	return request

}
