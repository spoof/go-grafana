// Copyright 2017 Sergey Safonov
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grafana

import (
	"encoding/json"
)

type Variables []Variable

// MarshalJSON implements json.Marshaler interface
func (v Variables) MarshalJSON() ([]byte, error) {
	vars := make([]probeVariable, len(v))
	for i, vv := range v {
		vars[i] = probeVariable{variable: vv}
	}

	jv := struct {
		List []probeVariable `json:"list"`
	}{
		List: vars,
	}

	return json.Marshal(jv)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (v *Variables) UnmarshalJSON(data []byte) error {
	jv := struct {
		List []probeVariable `json:"list"`
	}{}
	if err := json.Unmarshal(data, &jv); err != nil {
		return err
	}

	vars := make(Variables, len(jv.List))
	for i, v := range jv.List {
		vars[i] = v.variable
	}
	*v = vars

	return nil
}

type Variable interface {
	commonOptions() *commonVarOptions
}

type variableType string

const (
	intervalVarType   variableType = "interval"
	queryVarType      variableType = "query"
	datasourceVarType variableType = "datasource"
	customVarType     variableType = "custom"
	constantVarType   variableType = "constant"
)

type probeVariable struct {
	Type variableType `json:"type"`

	variable Variable
}

// MarshalJSON implements json.Marshaler interface
func (v *probeVariable) MarshalJSON() ([]byte, error) {
	// TODO: find a better way to implement this method
	var jj interface{}
	switch vv := v.variable.(type) {
	case *IntervalVariable:
		jj = struct {
			Type variableType `json:"type"`
			*IntervalVariable
		}{
			Type:             intervalVarType,
			IntervalVariable: vv,
		}
	case *QueryVariable:
		jj = struct {
			Type variableType `json:"type"`
			*QueryVariable
		}{
			Type:          queryVarType,
			QueryVariable: vv,
		}
	case *DatasourceVariable:
		jj = struct {
			Type variableType `json:"type"`
			*DatasourceVariable
		}{
			Type:               datasourceVarType,
			DatasourceVariable: vv,
		}
	case *CustomVariable:
		jj = struct {
			Type variableType `json:"type"`
			*CustomVariable
		}{
			Type:           customVarType,
			CustomVariable: vv,
		}
	case *ConstantVariable:
		jj = struct {
			Type variableType `json:"type"`
			*ConstantVariable
		}{
			Type:             constantVarType,
			ConstantVariable: vv,
		}
	}

	return json.Marshal(jj)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (v *probeVariable) UnmarshalJSON(data []byte) error {
	type JSONVariable probeVariable
	jv := struct {
		*JSONVariable
	}{
		JSONVariable: (*JSONVariable)(v),
	}
	if err := json.Unmarshal(data, &jv); err != nil {
		return err
	}

	var vv Variable
	switch jv.Type {
	// TODO handle unknown type
	case queryVarType:
		vv = new(QueryVariable)
	case intervalVarType:
		vv = new(IntervalVariable)
	case datasourceVarType:
		vv = new(DatasourceVariable)
	case customVarType:
		vv = new(CustomVariable)
	case constantVarType:
		vv = new(ConstantVariable)
	}
	if err := json.Unmarshal(data, vv); err != nil {
		return err
	}

	v.variable = vv

	return nil
}

type commonVarOptions struct {
	Name  string   `json:"name"`
	Label string   `json:"label"`
	Hide  hideType `json:"hide"`
}

type hideType uint

const (
	NoHide        hideType = 0
	HideLabelOnly hideType = 1
	HideVariable           = 2
)

// TODO: "refresh": 1, // it seems that On Dashboard Load / On Time Range change

type IntervalVariable struct {
	Auto      bool   `json:"auto"`
	StepCount uint   `json:"auto_count"`
	Min       uint   `json:"auto_min"`
	Query     string `json:"query"` // Values

	commonVarOptions
	/* TODO:
	"current": {
		"text": "1m",
		"value": "1m"
	},
	"options": [{
		"selected": false,
		"text": "auto",
		"value": "$__auto_interval"
	}]
	*/
}

func (v IntervalVariable) varType() variableType {
	return intervalVarType
}

func (v IntervalVariable) commonOptions() *commonVarOptions {
	return &v.commonVarOptions
}

type sortType uint

const (
	NoSort sortType = iota
	AlphabedicalASC
	AlphabedicalDESC
	NumericalASC
	NumericalDESC
)

type QueryVariable struct {
	Datasource string   `json:"datasource"`
	IncludeAll bool     `json:"includeAll"`
	Multi      bool     `json:"multi"`
	Query      string   `json:"query"`
	Regex      string   `json:"regex"`
	Sort       sortType `json:"sort"`
	AllValue   string   `json:"allValue"`

	/*
		TODO: "refresh": 2, // it seems that On Dashboard Load / On Time Range change
	*/
	commonVarOptions
}

func NewQueryVar(name string) *QueryVariable {
	return &QueryVariable{
		commonVarOptions: commonVarOptions{
			Name: name,
		},
	}
}

func (v *QueryVariable) commonOptions() *commonVarOptions {
	return &v.commonVarOptions
}

type DatasourceVariable struct {
	Query string `json:"query"` // it's datasource name
	Regex string `json:"regex"`

	commonVarOptions
}

func (v DatasourceVariable) commonOptions() *commonVarOptions {
	return &v.commonVarOptions
}

type CustomVariable struct {
	IncludeAll bool   `json:"includeAll"`
	AllValue   string `json:"allValue"`
	Multi      bool   `json:"multi"`
	Query      string `json:"query"` // comma separated value

	commonVarOptions
}

func (v CustomVariable) commonOptions() *commonVarOptions {
	return &v.commonVarOptions
}

// ConstantVariable is a dashboard variable of Constant type.
type ConstantVariable struct {
	Value string `json:"query"`
	commonVarOptions
}

// NewConstantVariable creates instance of CustomVariable with given name.
func NewConstantVariable(name string) *ConstantVariable {
	return &ConstantVariable{
		commonVarOptions: commonVarOptions{
			Name: name,
		},
	}
}

func (v ConstantVariable) commonOptions() *commonVarOptions {
	return &v.commonVarOptions
}
