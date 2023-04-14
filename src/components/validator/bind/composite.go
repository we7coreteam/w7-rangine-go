package bind

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

const defaultMemory = 32 << 20

type Composite struct {
	contentType string
}

func NewCompositeBind(contentType string) Composite {
	return Composite{
		contentType: contentType,
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

	switch c.contentType {
	case binding.MIMEJSON:
		return binding.JSON.Bind(req, obj)
	case binding.MIMEXML, binding.MIMEXML2:
		return binding.XML.Bind(req, obj)
	case binding.MIMEPROTOBUF:
		return binding.ProtoBuf.Bind(req, obj)
	case binding.MIMEMSGPACK, binding.MIMEMSGPACK2:
		return binding.MsgPack.Bind(req, obj)
	case binding.MIMEYAML:
		return binding.YAML.Bind(req, obj)
	case binding.MIMETOML:
		return binding.TOML.Bind(req, obj)
	default:
		return binding.Validator.ValidateStruct(obj)
	}
}
