// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package generator

import (
	"io"
)

type Generator interface {
	Generate(model interface{}, w io.Writer) error
}
