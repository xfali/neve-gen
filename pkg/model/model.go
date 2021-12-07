// Copyright (C) 2019-2021, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package model

type ModelData struct {
	Config Config
	Value  Value
}

func LoadModelData(configPath, valuePath string) (*ModelData, error) {
	c, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	v, err := LoadValue(valuePath)
	if err != nil {
		return nil, err
	}
	return &ModelData{
		Config: *c,
		Value:  *v,
	}, nil
}
