package common

import (
	"context"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/dialer"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/log"
	"git-biz.360es.cn/zion-infra/kamala-api/kamala/api"
	"git-biz.360es.cn/zion-infra/kamala-api/kamala/api/job"
)

//此文件用于kamala rpc客户端封装

func GreateJobByKamala(ctx context.Context, req *job.CreateJob_Request) error {
	conn, err := dialer.Dial(ctx, SERVICE_NAME_KAMALA)
	if err != nil {
		log.Errorf("[GreateJobByKamala]connect py rpc error:%v", err)
		return err

	}
	kamalaApi := api.NewJobClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()
	resJob, err := kamalaApi.CreateJob(ctx, req)
	log.Infof("[GreateJobByKamala]req :%+v  res is:%+v", req, resJob)
	if err != nil {
		log.Errorf("[GreateJobByKamala] create job  error:%v", err)
		return err

	}
	return nil

}
