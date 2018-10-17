package client

import (
	"github.com/spf13/viper"
)

type ClientInfo struct {
	srvInfo map[string]string
	cltInfo map[string]string
	srvName string // ServerName(eg:LobbyServer)
}

func GetNewClientInfo(srvName string) *ClientInfo {
	obj := &ClientInfo{srvName: srvName}
	obj.srvInfo = viper.GetStringMapString("ClientToServer")
	obj.cltInfo = viper.GetStringMapString("ServerToClient")
	return obj
}

func (ci *ClientInfo) GenWrap() (string, error) {
	var err error
	obj := GetNewCltWrap()
	obj.ReplaceSrvName(ci.srvName)
	if err = obj.ReplaceFunctions(ci.srvName, ci.srvInfo, ci.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, err
}

func (ci *ClientInfo) GenRegMsg() (string, error) {
	var err error
	obj := GetNewCltRegMsg()
	obj.ReplaceSrvName(ci.srvName)
	if err = obj.ReplaceFunctions(ci.srvName, ci.srvInfo, ci.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, err
}

func (ci *ClientInfo) GenExample() (string, error) {
	var err error
	obj := GetNewCltExample()
	obj.ReplaceSrvName(ci.srvName)
	if err = obj.ReplaceFunctions(ci.srvName, ci.srvInfo, ci.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, err
}
