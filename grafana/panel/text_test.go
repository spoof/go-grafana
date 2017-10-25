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

package panel

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	jsontools "github.com/spoof/go-grafana/pkg/json"
)

func TestTextPanel_MarshalJSON(t *testing.T) {
	panel := NewTextPanel(TextPanelMarkdownMode)
	panel.Content = "some content"

	got, err := json.MarshalIndent(panel, "", "\t\t")
	if err != nil {
		t.Fatalf("TextPanel.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"content": "some content",
		"mode": "markdown"
	}`)
	if eq, err := jsontools.BytesEqual(expected, got); err != nil {
		t.Fatalf("TextPanel.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("TextPanel.MarshalJSON: %s", pretty.Diff(expected, &got))
	}
}

func TestTextPanel_UnmarshalJSON(t *testing.T) {
	expected := NewTextPanel(TextPanelHTMLMode)
	expected.Content = "some content"

	data := []byte(`{
		"content": "some content",
		"mode": "html"
	}`)
	var got TextPanel
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("TextPanel.UnmarshalJSON returned error %s", err)
	}

	if !reflect.DeepEqual(expected, &got) {
		t.Errorf("TextPanel.UnmarshalJSON: %s", pretty.Diff(expected, &got))
	}
}
