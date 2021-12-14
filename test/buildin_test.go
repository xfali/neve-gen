// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-gen/pkg/buildin"
	"testing"
)

func TestWriteBuildin(t *testing.T) {
	err := buildin.WriteBuildinTemplate("../templates", "../pkg/buildin/buildin.go")
	if err != nil {
		t.Fatal(err)
	}
}
