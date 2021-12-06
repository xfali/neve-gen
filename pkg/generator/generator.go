// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package generator

import (
	"github.com/xfali/neve-gen/pkg/model"
	"io"
)

type Generator interface {
	Generate(model *model.TemplateModel, w io.Writer) error
}
