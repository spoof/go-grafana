package grafana

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTextPanel_MarshalJSON(t *testing.T) {
	panel := NewTextPanel(TextPanelMarkdownMode)
	panel.Content = "some content"
	options := panel.GeneralOptions()
	options.Title = "New Panel"
	options.Description = "Panel Description"
	options.Height = "250px"
	options.MinSpan = 1
	options.Span = 12
	options.Transparent = true

	got, err := json.MarshalIndent(panel, "", "\t\t")
	if err != nil {
		t.Fatalf("TextPanel.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"content": "some content",
		"mode": "markdown",
		"description": "Panel Description",
		"height": "250px",
		"links": null,
		"minSpan": 1,
		"span": 12,
		"title": "New Panel",
		"transparent": true,
		"id": 0,
		"type": "text"
	}`)
	if eq, err := JSONBytesEqual(expected, got); err != nil {
		t.Fatalf("TextPanel.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("TextPanel.MarshalJSON:\nexpected: %s\ngot: %s", expected, got)
	}
}

func TestTextPanel_UnmarshalJSON(t *testing.T) {
	expected := NewTextPanel(TextPanelHTMLMode)
	expected.Content = "some content"
	options := expected.GeneralOptions()
	options.Title = "New Panel"
	options.Description = "Panel Description"
	options.Height = "120px"
	options.MinSpan = 1
	options.Span = 12
	options.Transparent = true

	data := []byte(`{
		"content": "some content",
		"mode": "html",
		"description": "Panel Description",
		"height": "120px",
		"links": null,
		"minSpan": 1,
		"span": 12,
		"title": "New Panel",
		"transparent": true,
		"id": 0,
		"type": "text"
	}`)
	var got TextPanel
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("TextPanel.UnmarshalJSON returned error %s", err)
	}

	if !reflect.DeepEqual(expected, &got) {
		t.Errorf("TextPanel.MarshalJSON:\nexpected: %#v\ngot: %#v", expected, &got)
	}
}
