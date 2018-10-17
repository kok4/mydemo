package internal

//地图对象。如阻挡

import (
	"base/ape"
	"roomserver/conf"
	"roomserver/roommgr/room/internal/scn/internal/physic"
	"usercmd"
)

type _MapObject struct {
	apeObject *ape.RectangleParticle
	blockType usercmd.MapObjectConfigType
	typeId    int
}

func NewMapObject(typeid int, blockType usercmd.MapObjectConfigType, x, y, r float32) *_MapObject {
	return &_MapObject{
		apeObject: ape.NewRectangleParticle(x, y, r*2, r*2),
		blockType: blockType,
		typeId:    typeid}
}

// 类型: MapObjectConfigType_Block
func NewPhysicBlock(config *conf.MapNodeConfig) *_MapObject {
	mtype := usercmd.MapObjectConfigType(config.Type)
	blk := NewMapObject(config.Type, mtype, float32(config.Px), float32(config.Py), float32(config.Radius))
	blk.apeObject.SetFixed(true)
	return blk
}

func LoadMapObjectByConfig(config *conf.MapNodeConfig, physic *physic.ScenePhysic) {
	var obj *_MapObject
	switch usercmd.MapObjectConfigType(config.Type) {
	case usercmd.MapObjectConfigType_Block:
		obj = NewPhysicBlock(config)
	}
	physic.AddBlock(obj.apeObject)
}
