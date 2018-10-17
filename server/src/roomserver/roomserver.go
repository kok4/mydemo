package main

import (
	"base/glog"
	"fmt"
	"math/rand"
	"time"

	"roomserver/conf"
	"roomserver/gen/gensvr"
	"roomserver/roommgr/room"

	"github.com/spf13/viper"
)

type RoomServer struct {
	svr *gensvr.Server
}

func runRoomServer() {
	svr := RoomServer{}
	if svr.Init() {
		svr.Run()
	}
}

func (r *RoomServer) Init() bool {
	glog.Info("[启动] 开始初始化")

	rand.Seed(time.Now().Unix())

	// 全局配置
	if !conf.ConfigMgr_GetMe().Init() {
		glog.Error("[启动] 读取全局配置失败")
		return false
	}

	// 绑定本地端口
	address := viper.GetString("room.listen")
	svr, err := gensvr.New("tcp", address, 10000)
	if err != nil {
		glog.Error(fmt.Sprintf("创建服务器失败: ", err))
		return false
	}
	r.svr = svr

	if !room.Init() {
		return false
	}
	if !conf.InitMapConfig() {
		glog.Error("[启动]InitMapConfig fail! ")
		return false
	}

	glog.Info("[启动] 完成初始化")
	return true
}

func (r *RoomServer) Run() {
	r.svr.Run()
}
