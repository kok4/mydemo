package conf

import (
	"base/glog"
	"encoding/xml"
	"io/ioutil"
	"roomserver/types"
	"strconv"
	"strings"
	"sync"
	"usercmd"

	"github.com/spf13/viper"
)

//刷新点
type XmlFoodPoint struct {
	ObjType int32                 `xml:"objtype,attr"`
	Foods   []XmlFoodPointRefresh `xml:"food"`
}

type XmlFoodPointRefresh struct {
	ID       uint16  `xml:"id,attr"`
	Num      uint32  `xml:"num,attr"`
	Interval float64 `xml:"interval,attr"`
}

type XmlScene struct {
	Id types.SceneID `xml:"id,attr"`
}

type XmlMap struct {
	XMLName xml.Name   `xml:"config"`
	Scenes  []XmlScene `xml:"scene"`
}

type XmlFoodItem struct {
	FoodId    uint16                `xml:"id,attr"`
	FoodType  uint16                `xml:"type,attr"`
	Size      float32               `xml:"size,attr"`
	BirthTime float64               `xml:"birthTime,attr"`
	MapNum    uint16                `xml:"mapnum,attr"`
	HP        uint32                `xml:"hp,attr"`
	Buffstr   string                `xml:"buff,attr"`
	Area      int                   `xml:"area,attr"`
	Exp       uint32                `xml:"exp,attr"`
	LiveTime  int64                 `xml:"time,attr"`
	Children  []XmlFoodPointRefresh `xml:"child"`
	Buff      []uint32
	Rate      []uint16 // 概率
	Sum       uint16   // 基数
}

type XmlFoodItems struct {
	MapId  types.SceneID `xml:"mapid,attr"`
	ShotID uint16        `xml:"shotid,attr"`
	Items  []XmlFoodItem `xml:"item"`
}

type XmlFoodCfg struct {
	XMLName xml.Name       `xml:"config"`
	Foods   []XmlFoodItems `xml:"food"`
}

//////////////////////////////////////////////////////////////////////////////////
//配置类
type ConfigMgr struct {
	Map       *XmlMap
	FoodDatas *XmlFoodCfg
}

var (
	configm      *ConfigMgr
	configmMutex sync.RWMutex
)

func NewConfigMgr() *ConfigMgr {
	return &ConfigMgr{}
}

func ConfigMgr_GetMe() (c *ConfigMgr) {
	if configm == nil {
		configm = NewConfigMgr()
	}
	configmMutex.RLock()
	c = configm
	configmMutex.RUnlock()
	return
}

func ReloadConfig() bool {
	c := NewConfigMgr()
	if !c.Init() {
		return false
	}
	configmMutex.Lock()
	configm = c
	configmMutex.Unlock()
	return true
}

// 全局map
func (r *ConfigMgr) LoadMap() bool {
	content, err := ioutil.ReadFile(viper.GetString("global.xmlcfg") + "map.xml")
	if err != nil {
		glog.Error("[配置] 打开map配置失败 ", err)
		return false
	}
	xmlmap := &XmlMap{}
	err = xml.Unmarshal(content, xmlmap)
	if err != nil {
		glog.Error("[配置] 解析map配置失败 ", err)
		return false
	}
	r.Map = xmlmap
	glog.Info("LoadMap:", len(r.Map.Scenes))
	return true
}

func (r *ConfigMgr) LoadFoodCfg() bool {
	content, err := ioutil.ReadFile(viper.GetString("global.xmlcfg") + "food.xml")
	if err != nil {
		glog.Error("[配置] 打开配置 food.xml失败 ", err)
		return false
	}
	foodData := &XmlFoodCfg{}
	err = xml.Unmarshal(content, foodData)
	if err != nil {
		glog.Error("[配置] 解析配置 food.xml 失败", err)
		return false
	}
	for index, food := range foodData.Foods {
		for innr, item := range food.Items {
			if len(item.Buffstr) == 0 {
				continue
			}
			strs := strings.Split(item.Buffstr, "|")
			for _, v := range strs {
				tmp := strings.Split(v, ":")
				if len(tmp) != 2 {
					glog.Error("food.xml bufferandrate error", v)
					return false
				}
				buff, ok := strconv.Atoi(tmp[0])
				if nil != ok {
					glog.Error("food.xml buff error: ", strs, ",", item.Buffstr)
					return false
				}
				rate, ok := strconv.Atoi(tmp[1])
				foodData.Foods[index].Items[innr].Buff = append(foodData.Foods[index].Items[innr].Buff, uint32(buff))
				foodData.Foods[index].Items[innr].Rate = append(foodData.Foods[index].Items[innr].Rate, uint16(rate))
				foodData.Foods[index].Items[innr].Sum += uint16(rate)
			}
			if 0 == foodData.Foods[index].Items[innr].Sum ||
				len(foodData.Foods[index].Items[innr].Buff) != len(foodData.Foods[index].Items[innr].Rate) {
				glog.Error("food.xml sum is zero or len error", strs)
				return false
			}
		}
	}
	r.FoodDatas = foodData
	return true
}

func (r *ConfigMgr) Init() bool {
	if !r.LoadMap() {
		return false
	}

	if !r.LoadFoodCfg() {
		return false
	}

	glog.Info("[配置] 加载配置成功 ")
	return true
}

func (r *ConfigMgr) GetXmlFoodItems(sceneId types.SceneID) *XmlFoodItems {
	for _, m := range r.FoodDatas.Foods {
		if m.MapId == sceneId {
			return &m
		}
	}
	return nil
}

//food 表
func (r *ConfigMgr) GetFood(sceneID types.SceneID, foodid uint16) *XmlFoodItem {
	item := r.GetXmlFoodItems(sceneID)
	for _, val := range item.Items {
		if val.FoodId == foodid {
			return &val
		}
	}
	return nil
}

func (r *ConfigMgr) GetFoodSize(sceneID types.SceneID, typeId uint16) float32 {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return food.Size
	}
	return 0.5
}

func (r *ConfigMgr) GetFoodMapNum(sceneID types.SceneID, typeId uint16) uint16 {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return food.MapNum
	}
	return 0
}

func (r *ConfigMgr) GetFoodHP(sceneID types.SceneID, typeId uint16) uint32 {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return food.HP
	}
	return 0
}

func (r *ConfigMgr) GetFoodExp(sceneID types.SceneID, typeId uint16) uint32 {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return food.Exp
	}
	return 0
}

func (r *ConfigMgr) GetFoodTime(sceneID types.SceneID, typeId uint16) int64 {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return food.LiveTime
	}
	return 0
}

func (r *ConfigMgr) GetFoodBallType(sceneID types.SceneID, typeId uint16) usercmd.BallType {
	if food := r.GetFood(sceneID, typeId); food != nil {
		return usercmd.BallType(food.FoodType)
	}
	return 0
}
