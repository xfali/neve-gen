// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"context"
	"github.com/xfali/neve-gen/pkg/model"
)

type Stage interface {
	Name() string
	Generate(ctx context.Context, model *model.ModelData) error
	Rollback(ctx context.Context) error
}
