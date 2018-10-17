package server

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServerInfo struct {
	srvInfo map[string]string
	cltInfo map[string]string
	srvName string // ServerName(eg:LobbyServer)
}

func GetNewServerInfo(srvName string) *ServerInfo {
	obj := &ServerInfo{srvName: srvName}
	obj.srvInfo = viper.GetStringMapString("ClientToServer")
	obj.cltInfo = viper.GetStringMapString("ServerToClient")
	return obj
}

func (srv *ServerInfo) GenWrap() (string, error) {
	var err error
	obj := GetNewSrvWrap()
	obj.ReplaceSrvName(srv.srvName)
	if err = obj.ReplaceFunctions(srv.srvName, srv.srvInfo, srv.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, nil
}

func (srv *ServerInfo) GenRegMsg() (string, error) {
	var err error
	obj := GetNewSrvRegMsg()
	obj.ReplaceSrvName(srv.srvName)
	if err := obj.ReplaceFunctions(srv.srvName, srv.srvInfo, srv.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, err
}

func (srv *ServerInfo) GenServer() (string, error) {
	//var err error
	//obj := GetNewGenServer()
	//obj.ReplaceSrvName(srv.srvName)
	//if err := obj.ReplaceFunctions(srv.srvName, srv.srvInfo, srv.cltInfo); err != nil {
	//	return "", err
	//}
	//str := obj.GetContent()
	//fmt.Println(str)
	//return str, err
	return "", fmt.Errorf("Interface has been abandoned. ")
}

func (srv *ServerInfo) GenExample() (string, error) {
	var err error
	obj := GetNewSrvExample()
	obj.ReplaceSrvName(srv.srvName)
	if err := obj.ReplaceFunctions(srv.srvName, srv.srvInfo, srv.cltInfo); err != nil {
		return "", err
	}
	str := obj.GetContent()
	//fmt.Println(str)
	return str, err
}
