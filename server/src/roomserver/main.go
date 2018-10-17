// Room server.
package main

import (
	"base/glog"
	"flag"
	"fmt"
	"runtime/debug"

	"github.com/spf13/viper"
)

var (
	config = flag.String("config", "config.json", "config path")
)

func main() {
	flag.Parse()
	defer func() {
		if err := recover(); err != nil {
			glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
		}
		glog.Flush()
	}()

	viper.SetConfigFile(*config)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %s", err))
	}

	initLog()
	runRoomServer()
	glog.Info("[关闭] 房间服务器关闭完成")
} // main()

func initLog() {
	loglevel := viper.GetString("global.loglevel")
	if loglevel != "" {
		flag.Lookup("stderrthreshold").Value.Set(loglevel)
	}

	glog.SetLogFile(viper.GetString("room.log"))
	glog.SetLogType(glog.LogTimeType_Day) //日志每天一换
} // initLog()
