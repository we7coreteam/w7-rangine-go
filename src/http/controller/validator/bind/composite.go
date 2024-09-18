package bind

import (
	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"reflect"
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

	//预先判断一次是context.shouldBind方法里面得multipart解析使用的包大小设置是写死的
	if c.context.ContentType() == binding.MIMEMultipartPOSTForm {
		_, err := c.context.MultipartForm()
		if err != nil {
			return err
		}
	}

	var err error
	b := binding.Default(req.Method, c.context.ContentType())
	if reflect.TypeOf(b).Implements(reflect.TypeOf((*binding.BindingBody)(nil)).Elem()) {
		err = c.context.ShouldBindBodyWith(obj, b.(binding.BindingBody))
	} else {
		err = c.context.ShouldBindWith(obj, b)
	}
	if err == nil {
		err = defaults.Set(obj)
	}

	return err
}
