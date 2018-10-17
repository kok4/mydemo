package roommgr

// 房间管理类

import (
	"runtime/debug"
	"sync"
	"time"

	"base/glog"
	rm "roomserver/roommgr/room"
	"roomserver/types"
)

type RoomMgr struct {
	mutex sync.RWMutex // protect rooms
	rooms map[types.RoomID]*rm.Room
}

var (
	roommgr      *RoomMgr
	roommgr_once sync.Once
)

func GetMe() *RoomMgr {
	if roommgr == nil {
		roommgr_once.Do(func() {
			roommgr = &RoomMgr{
				rooms: make(map[types.RoomID]*rm.Room),
			}
			roommgr.init()
		})
	}
	return roommgr
}

func (r *RoomMgr) init() {
	go func() {
		tick := time.NewTicker(time.Second * 1)
		defer func() {
			if err := recover(); err != nil {
				glog.Error("[异常] 报错 ", err, "\n", string(debug.Stack()))
			}
			tick.Stop()
		}()
		for {
			select {
			case <-tick.C:
				r.timeAction()
			}
		}
	}()
}

// 定时事件
func (r *RoomMgr) timeAction() {
	r.mutex.RLock()
	rms := r.getClosedRooms()
	r.mutex.RUnlock()

	// 删除已关闭房间.
	for _, room := range rms {
		r.RemoveRoom(room.ID())
	}
}

// 添加房间
func (r *RoomMgr) addRoom(room *rm.Room) *rm.Room {
	r.rooms[room.ID()] = room
	return room
}

// 删除房间
func (r *RoomMgr) RemoveRoom(rid types.RoomID) {
	r.mutex.Lock()
	delete(r.rooms, rid)
	r.mutex.Unlock()
}

func (r *RoomMgr) NewRoom(sceneId types.SceneID) *rm.Room {
	r.mutex.Lock()
	room := r.addRoom(rm.New(sceneId))
	r.mutex.Unlock()

	glog.Infof("[房间] 创建: sceneID=%d roomID=%d", sceneId, room.ID())
	return room
}

func (r *RoomMgr) GetRoomByID(roomID types.RoomID) *rm.Room {
	if roomID == 0 {
		return nil // 0 为非法ID
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, ok := r.rooms[roomID]
	if !ok {
		return nil
	}
	return room
}

// getClosedRooms 生成已关闭房间列表.
func (r *RoomMgr) getClosedRooms() []*rm.Room {
	rms := make([]*rm.Room, 0)
	for _, room := range r.rooms {
		if room.IsClosed() {
			rms = append(rms, room)
		}
	}
	return rms
}
