// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package main

import (
	"flag"
	"github.com/xfali/neve-gen/pkg/layout"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/neve-gen/pkg/stage"
	"github.com/xfali/xlog"
)

func main() {
	target := flag.String("o", "./awesomeProject", "Target directory")
	tmplPath := flag.String("t", "", "Template directory")
	confPath := flag.String("c", "", "Config file path")
	valuePath := flag.String("f", "", "Value file path")
	flag.Parse()
	if *target == "" {
		xlog.Fatalln("Target directory is empty, set it with flag -o")
	}
	if *confPath == "" {
		xlog.Fatalln("Config file path is empty, set it with flag -c")
	}
	if *valuePath == "" {
		xlog.Fatalln("Value file path is empty, set it with flag -f")
	}

	var stages []stage.Stage
	var err error
	if *tmplPath == "" {
		xlog.Warnln("Template directory is empty, use buildin template. You can set it with flag -t")
		stages, err = layout.ParseBuildinStages(*target)
		if err != nil {
			xlog.Fatalln(err)
		}
	} else {
		stages, err = layout.ParseStages(*target, *tmplPath)
		if err != nil {
			xlog.Fatalln(err)
		}
	}

	app := layout.NewProjectGenerator(stages)
	m ,err := model.LoadModelData(*confPath, *valuePath)
	if err != nil {
		xlog.Fatalln(err)
	}
	err = app.Layout(m)
	if err != nil {
		xlog.Fatalln(err)
	}
}
