package envconfig

import (
	"github.com/ErickMaria/envcontainer/internal/pkg/handler/errors"
	"github.com/ErickMaria/envcontainer/internal/pkg/syscmd"
)

type Template struct {
	PathDefault string
}

func newTemplate() Template {

	return Template{
		PathDefault: Home(),
	}
}

func CreateIfNotExists() {
	template := newTemplate()
	if !syscmd.ExistsPath(template.PathDefault) {
		template.Init()
	}
}

func (template Template) Init() {
	err := syscmd.CreatePath(template.PathDefault)
	errors.Throw("envcontainer: error to create folders, check permissions", err)
}
