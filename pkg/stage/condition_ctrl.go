// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package stage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/xfali/neve-gen/pkg/database"
	"github.com/xfali/neve-gen/pkg/model"
	"github.com/xfali/xlog"
	"strings"
	"text/template"
)

func CheckCondition(ctx context.Context, condition string, m *model.ModelData) bool {
	if condition != "" {
		if strings.Index(condition, "${MODULE}") != -1 {
			return true
		}
		tmp, err := template.New("checkCondition").Funcs(map[string]interface{}{
			"loadFromDB": func(name string) bool { return database.IsTable(ctx, name) },
			"notFromDB":  func(name string) bool { return !database.IsTable(ctx, name) },
		}).Parse(fmt.Sprintf("{{if %s}}1{{else}}0{{end}}", condition))
		if err != nil {
			xlog.Warnf("Condition [%s] parse error: %v\n", condition, err)
			return true
		}
		b := bytes.Buffer{}
		b.Grow(1)
		err = tmp.Execute(&b, m)
		if err != nil {
			xlog.Warnf("Condition [%s] parse error: %v\n", condition, err)
			return true
		}
		xlog.Debugln("Condition: ", condition, " ", b.String())
		return b.String() == "1"
	}
	return true
}
