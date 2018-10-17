package birth

import (
	"roomserver/types"
)

type IScene interface {
	GetRandPos() (x, y float64)
	SceneID() types.SceneID
	GenBallID() uint32
}
