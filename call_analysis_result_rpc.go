package common

import (
	"context"
	"errors"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/dialer"
	"git-biz.360es.cn/infra-components/sdk/microservice-framework/go-framework.git/log"
	"zmind-service/interal_api"
)

// 此文件用于python rpc客户端封装
//得到 cocuckoo沙箱的解析结果

func CallGetAnalysisRestultByRpc(ctx context.Context, logPath, assetName string) (string, error) {

	conn, err := dialer.Dial(ctx, SERVICE_NAME_INTERAL_RPC, dialer.WithRegistry(nil))
	if err != nil {
		log.Errorf("[CallGetCuckooRestByPyRpc]connect py rpc error:%v", err)
		return "", err

	}
	req := new(interal_api.Request)
	req.AssetName = assetName
	req.LogPath = logPath
	pyApi := interal_api.NewNotifyClient(conn)
	//nc, cancel := context.WithCancel(ctx)
	//defer cancel()

	// 调用pythonrpc服务
	rpcRes, err := pyApi.GetAnalysisResult(ctx, req)
	if err != nil {
		log.Errorf("[CallGetCuckooRestByPyRpc]call py error:%v", err)
		return "", err
	}
	if rpcRes.Code != interal_api.Code_OK {
		log.Infof("GetCuckooRes return error:%v", rpcRes.Info)
		return "", errors.New(rpcRes.Info)
	}

	log.Infof("response %v", rpcRes)
	return rpcRes.Result, nil

}
