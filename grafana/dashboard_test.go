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

	"github.com/kr/pretty"
	"github.com/spoof/go-grafana/grafana/panel"
	"github.com/spoof/go-grafana/pkg/field"
)

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func TestDashboard_MarshalJSON(t *testing.T) {
	d := NewDashboard("Dashboard Title")
	d.Editable = true
	d.GraphTooltip = 2
	d.HideControls = true
	d.Style = dashboardLightStyle
	d.Timezone = "MSK"
	d.tags = field.NewTags("tag1", "tag2")

	row1 := NewRow()
	row1.Collapsed = true
	row1.Height = field.ForceString("200")
	row1.RepeatFor = "variable1"
	row1.ShowTitle = true
	row1.Title = "Row Title 1"
	row1.TitleSize = "h1"
	p1 := panel.NewText(panel.TextPanelHTMLMode)
	p1.Content = "Content 1"
	opts1 := p1.GeneralOptions()
	opts1.Description = "Panel Description 1"
	opts1.Height = "200px"
	opts1.MinSpan = 1
	opts1.Span = 12
	opts1.Title = "Panel Title 1"
	opts1.Transparent = true
	row1.Panels = []Panel{p1}

	row2 := NewRow()
	row2.Collapsed = true
	row2.Height = field.ForceString("200")
	row2.RepeatFor = "variable2"
	row2.ShowTitle = true
	row2.Title = "Row Title 2"
	row2.TitleSize = "h2"
	p2 := panel.NewText(panel.TextPanelHTMLMode)
	p2.Content = "Content 2"
	opts2 := p2.GeneralOptions()
	opts2.Description = "Panel Description 2"
	opts2.Height = "200px"
	opts2.MinSpan = 1
	opts2.Span = 12
	opts2.Title = "Panel Title 2"
	opts2.Transparent = true
	row2.Panels = []Panel{p2}

	d.Rows = []*Row{row1, row2}

	got, err := json.MarshalIndent(&d, "", "\t")
	if err != nil {
		t.Fatalf("Dashboard.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"schemaVersion": 14,
		"editable": true,
		"graphTooltip": 2,
		"hideControls": true,
		"rows": [{
			"collapse": true,
			"editable": true,
			"height": "200",
			"panels": [{
				"id": 1,
				"type": "text",
				"description": "Panel Description 1",
				"height": "200px",
				"links": null,
				"minSpan": 1,
				"span": 12,
				"title": "Panel Title 1",
				"transparent": true,

				"content": "Content 1",
				"mode": "html"
			}],
			"repeat": "variable1",
			"showTitle": true,
			"title": "Row Title 1",
			"titleSize": "h1"
		},
		{
			"collapse": true,
			"editable": true,
			"height": "200",
			"panels": [{
				"id": 2,
				"type": "text",
				"description": "Panel Description 2",
				"height": "200px",
				"links": null,
				"minSpan": 1,
				"span": 12,
				"title": "Panel Title 2",
				"transparent": true,

				"content": "Content 2",
				"mode": "html"
			}],
			"repeat": "variable2",
			"showTitle": true,
			"title": "Row Title 2",
			"titleSize": "h2"
		}],
		"style": "light",
		"timezone": "MSK",
		"title": "Dashboard Title",
		"tags": ["tag1", "tag2"]
	}`)
	if eq, err := JSONBytesEqual(expected, got); err != nil {
		t.Fatalf("Dashboard.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("Dashboard.MarshalJSON: got %s, want %s\n", got, expected)
	}
}

func TestDashboard_UnmarshalJSON(t *testing.T) {
	// TODO: add meta unmarshaling
	data := []byte(`{
		"id": 1,
		"version": 2,
		"editable": true,
		"graphTooltip": 2,
		"hideControls": true,
		"style": "light",
		"timezone": "msk",
		"title": "Dashboard Title",
		"tags": ["tag1", "tag2"],
		"schemaVersion": 12
	}`)
	var got Dashboard
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Dashboard.UnmarshalJSON returned error %s", err)
	}

	expected := NewDashboard("Dashboard Title")
	expected.ID = 1
	expected.SchemaVersion = 12
	expected.Version = 2
	expected.Editable = true
	expected.GraphTooltip = 2
	expected.HideControls = true
	expected.Style = dashboardLightStyle
	expected.Timezone = "msk"
	expected.tags = field.NewTags("tag1", "tag2")

	if !reflect.DeepEqual(&got, expected) {
		t.Errorf("Dashboard.UnmarshalJSON: %s", pretty.Diff(&got, expected))
	}

}

func TestProbePanel_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"description": "Panel Description",
		"height": "250px",
		"links": null,
		"minSpan": 1,
		"span": 12,
		"title": "New Panel",
		"transparent": true,
		"id": 1,
		"type": "text",
		"mode": "text",
		"content": "Content"
	}`)
	var got probePanel
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("probePanel.UnmarshalJSON returned error %s", err)
	}

	if got.GeneralOptions() == nil {
		t.Errorf("probePanel.UnmarshalJSON. got.GeneralOptions() shouldn't be nil")
	}

	epxectedPanel := panel.NewText(panel.TextPanelMarkdownMode)
	epxectedPanel.Content = "Content"
	epxectedPanel.Mode = panel.TextPanelTextMode
	expected := &probePanel{ID: 1, Type: textPanel, panel: epxectedPanel}
	opts := expected.GeneralOptions()
	opts.Description = "Panel Description"
	opts.Height = "250px"
	opts.MinSpan = 1
	opts.Span = 12
	opts.Title = "New Panel"
	opts.Transparent = true
	if !reflect.DeepEqual(&got, expected) {
		t.Errorf("probePanel.UnmarshalJSON. got: %+v, want: %+v", got, expected)
	}
}

func TestProbePanel_MarshalJSON(t *testing.T) {
	panel := panel.NewText(panel.TextPanelMarkdownMode)
	panel.Content = "Content"
	opts := panel.GeneralOptions()
	opts.Description = "Panel Description"
	opts.Title = "Singlestat Panel"
	opts.Height = "250px"
	opts.MinSpan = 1
	opts.Span = 12
	opts.Transparent = true
	pp := &probePanel{ID: 1, Type: textPanel, panel: panel}

	got, err := json.MarshalIndent(pp, "", "\t\t")
	if err != nil {
		t.Fatalf("probePanel.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"description": "Panel Description",
		"height": "250px",
		"links": null,
		"minSpan": 1,
		"span": 12,
		"title": "Singlestat Panel",
		"transparent": true,
		"id": 1,
		"type": "text",

		"content": "Content",
		"mode": "markdown"
	}`)
	if eq, err := JSONBytesEqual(expected, got); err != nil {
		t.Fatalf("probePanel.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("probePanel.MarshalJSON: got %s, want %s\n", got, expected)
	}
}
