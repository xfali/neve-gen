// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package layout

import "github.com/xfali/neve-gen/pkg/stage"

func CreateStages(target, tmplRoot string) []stage.Stage {
	var ret []stage.Stage
	ret = append(ret, stage.NewGenProjectStage(target, tmplRoot))
	ret = append(ret, stage.NewGenGoModStage(target, tmplRoot))
	ret = append(ret, stage.NewGenMainStage(target, tmplRoot))
	ret = append(ret, stage.NewGenHandlerStage(target, tmplRoot))
	ret = append(ret, stage.NewGenServiceStage(target, tmplRoot))
	ret = append(ret, stage.NeGenServiceImplStage(target, tmplRoot))
	return ret
}
