// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/xfali/neve-gen/pkg/layout"
	"github.com/xfali/neve-gen/pkg/model"
	"testing"
)

func TestAppLayout(t *testing.T) {
	stages, err := layout.ParseStages("../../testproj", "../templates")
	if err != nil {
		t.Fatal(err)
	}
	app := layout.NewProjectGenerator(stages)
	m ,err := model.LoadModelData("test-conf.yaml", "test-value.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = app.Layout(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAppBuildinLayout(t *testing.T) {
	stages, err := layout.ParseBuildinStages("../../testproj")
	if err != nil {
		t.Fatal(err)
	}
	app := layout.NewProjectGenerator(stages)
	m ,err := model.LoadModelData("test-conf.yaml", "test-value.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = app.Layout(m)
	if err != nil {
		t.Fatal(err)
	}
}
