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
	"reflect"
	"testing"
)

func TestVariables_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"list": [{
			"name": "var",
			"label": "Variable",
			"hide": 2,
			"type": "query",

			"datasource": "Prometheus",
			"includeAll": true,
			"multi": true,
			"query": "up{job=\"prometheus\"}",
			"regex": "/local/",
			"sort": 4,
			"allValue": ".*"
		}]
	}`)
	var got Variables
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Variables.UnmarshalJSON returned error %s", err)
	}

	v := NewQueryVar("var")
	v.Label = "Variable"
	v.Hide = HideVariable
	v.Query = `up{job="prometheus"}`
	v.IncludeAll = true
	v.Multi = true
	v.Datasource = "Prometheus"
	v.AllValue = ".*"
	v.Regex = "/local/"
	v.Sort = NumericalDESC
	expected := Variables{v}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Variables.UnmarshalJSON. got: %+v, want: %+v", got, expected)
	}
}

func TestVariables_MarshalJSON(t *testing.T) {
	v := NewQueryVar("var")
	v.Label = "Variable"
	v.Hide = HideVariable
	v.Query = `up{job="prometheus"}`
	v.IncludeAll = true
	v.Multi = true
	v.Datasource = "Prometheus"
	v.AllValue = ".*"
	v.Regex = "/local/"
	v.Sort = NumericalDESC
	variables := Variables{v}

	got, err := json.MarshalIndent(variables, "", "\t\t")
	if err != nil {
		t.Fatalf("Variables.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"list": [{
			"name": "var",
			"label": "Variable",
			"hide": 2,
			"type": "query",
			"datasource": "Prometheus",
			"includeAll": true,
			"multi": true,
			"query": "up{job=\"prometheus\"}",
			"regex": "/local/",
			"sort": 4,
			"allValue": ".*"
		}]
	}`)
	if eq, err := JSONBytesEqual(expected, got); err != nil {
		t.Fatalf("probeVariable.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("probeVariable.MarshalJSON: got %s, want %s\n", got, expected)
	}
}
