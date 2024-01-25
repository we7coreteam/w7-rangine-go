package bind

import (
	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Composite struct {
	context *gin.Context
}

func NewCompositeBind(ctx *gin.Context) Composite {
	return Composite{
		context: ctx,
	}
}

func (Composite) Name() string {
	return "composite"
}

func (c Composite) Bind(req *http.Request, obj any) error {
	if c.context.Params != nil {
		m := make(map[string][]string)
		for _, v := range c.context.Params {
			m[v.Key] = []string{v.Value}
		}
		err := binding.MapFormWithTag(obj, m, "uri")
		if err != nil {
			return err
		}
	}

	err := c.context.ShouldBind(obj)
	if err == nil {
		err = defaults.Set(obj)
	}

	return err
}
