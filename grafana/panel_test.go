package grafana

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPanelHeight_UnmarshalJSON(t *testing.T) {
	tt := []struct {
		got      string
		expected panelHeight
	}{
		{got: `{"height": ""}`, expected: ""},
		{got: `{"height": null}`, expected: ""},
		{got: `{"height": 200}`, expected: "200"},
		{got: `{"height": "200px"}`, expected: "200px"},
	}

	for _, ts := range tt {
		data := struct {
			Height panelHeight `json:"height"`
		}{}

		if err := json.Unmarshal([]byte(ts.got), &data); err != nil {
			t.Fatalf("panelHeight.UnmarshalJSON returned error %s", err)
		}

		if !reflect.DeepEqual(ts.expected, data.Height) {
			t.Errorf("panelHeight.MarshalJSON:\nexpected: %#v\ngot: %#v", ts.expected, data.Height)
		}
	}
}
