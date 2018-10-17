# 目录说明

## 子目录

* `conf`		配置包，被其他包依赖
* `match`		匹配功能，依赖 roommgr 和 sess
* `msgdef`      消息定义和生成代码, 不依赖其他包
* `roommgr`		房间管理，内含房间及房间内游戏逻辑
* `sess`		玩家会话(session)，依赖 roommgr
* `types`		公共类型定义

## 文件

* `main.go`  			主程序
* `READMD.md`			本文件
* `roomserver.go`		房间服务器
* `wilds_MsgProc.go`	wilds消息处理器
