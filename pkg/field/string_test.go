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

package field

import (
	"encoding/json"
	"reflect"
	"testing"

	jsontools "github.com/spoof/go-grafana/pkg/json"
)

func TestForceString_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		got      string
		expected ForceString
	}{
		{got: `{"height": ""}`, expected: ""},
		{got: `{"height": null}`, expected: ""},
		{got: `{"height": 200}`, expected: "200"},
		{got: `{"height": "200px"}`, expected: "200px"},
	}

	for _, ts := range tt {
		data := struct {
			Height ForceString `json:"height"`
		}{}

		if err := json.Unmarshal([]byte(ts.got), &data); err != nil {
			t.Fatalf("ForceString.UnmarshalJSON returned error %s", err)
		}

		if !reflect.DeepEqual(ts.expected, data.Height) {
			t.Errorf("ForceString.MarshalJSON:\nexpected: %#v\ngot: %#v", ts.expected, data.Height)
		}
	}
}

func TestForceString_MarshalJSON(t *testing.T) {
	data := struct {
		Height ForceString `json:"height"`
	}{
		Height: ForceString("200"),
	}

	got, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("ForceString.MarshalJSON returned error %s", err)
	}

	expected := []byte(`{"height": "200"}`)
	if eq, err := jsontools.BytesEqual(expected, got); err != nil {
		t.Fatalf("ForceString.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("ForceString.MarshalJSON: got %s, want %s\n", got, expected)
	}
}
