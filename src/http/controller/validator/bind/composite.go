package bind

import (
	"errors"
	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

const defaultMemory = 32 << 20

type Composite struct {
	contentType string
	uriParams   gin.Params
}

func NewCompositeBind(contentType string, uriParams gin.Params) Composite {
	return Composite{
		contentType: contentType,
		uriParams:   uriParams,
	}
}

func (Composite) Name() string {
	return "composite"
}

func (c Composite) Bind(req *http.Request, obj any) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
		return err
	}
	if err := binding.MapFormWithTag(obj, req.Form, "form"); err != nil {
		return err
	}
	if c.uriParams != nil {
		m := make(map[string][]string)
		for _, v := range c.uriParams {
			m[v.Key] = []string{v.Value}
		}
		if err := binding.MapFormWithTag(obj, m, "uri"); err != nil {
			return err
		}
	}

	var err error
	switch c.contentType {
	case binding.MIMEJSON:
		err = binding.JSON.Bind(req, obj)
	case binding.MIMEXML, binding.MIMEXML2:
		err = binding.XML.Bind(req, obj)
	case binding.MIMEPROTOBUF:
		err = binding.ProtoBuf.Bind(req, obj)
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		err = binding.MsgPack.Bind(req, obj)
	case binding.MIMEYAML:
		err = binding.YAML.Bind(req, obj)
	case binding.MIMETOML:
		err = binding.TOML.Bind(req, obj)
	default:
		err = binding.Validator.ValidateStruct(obj)
	}

	if err == nil {
		err = defaults.Set(obj)
	}

	return err
}
