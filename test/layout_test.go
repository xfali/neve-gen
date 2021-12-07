// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-gen/pkg/layout"
	"testing"
)

func TestAppLayout(t *testing.T) {
	stages, err := layout.ParseStages("./proj", "../template")
	if err != nil {
		t.Fatal(err)
	}
	app := layout.NewProjectGenerator(stages)
	err = app.Layout("./proj", getTestModel(t))
	if err != nil {
		t.Fatal(err)
	}
}
