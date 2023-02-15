package global

import (
	"github.com/we7coreteam/w7-rangine-go/src/app"
)

type RGlobal struct {
	App *app.App
}

var G = new(RGlobal)
