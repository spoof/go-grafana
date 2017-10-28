package panel

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kr/pretty"
	"github.com/spoof/go-grafana/pkg/field"
	jsontools "github.com/spoof/go-grafana/pkg/json"
)

func TestGraph_MarshalJSON(t *testing.T) {
	panel := NewGraph()
	leftMin := field.ForceString("0")
	rightMax := field.ForceString("10")

	panel.YAxes = GraphYaxesOptions{
		Left: GraphYAxis{
			Format:  "dtdurations",
			Label:   "label1",
			LogBase: 1,
			Max:     nil,
			Min:     &leftMin,
			Show:    true,
		},
		Right: GraphYAxis{
			Format:  "short",
			Label:   "",
			LogBase: 10,
			Max:     &rightMax,
			Min:     nil,
			Show:    false,
		},
	}
	got, err := json.MarshalIndent(panel, "", "\t")
	if err != nil {
		t.Fatalf("Graph.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"yaxes": [{
			"format": "dtdurations",
			"label": "label1",
			"logBase": 1,
			"max": null,
			"min": "0",
			"show": true
		},
		{
			"format": "short",
			"logBase": 10,
			"max": "10",
			"min": null,
			"show": false
		}]
	}`)
	if eq, err := jsontools.BytesEqual(expected, got); err != nil {
		t.Fatalf("Graph.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("Graph.MarshalJSON:\ngot %s\nwant: %s", expected, got)
	}
}
func TestGraph_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"yaxes": [{
			"format": "dtdurations",
			"label": "label1",
			"logBase": 1,
			"max": null,
			"min": "0",
			"show": true
		},
		{
			"format": "short",
			"logBase": 10,
			"max": null,
			"min": null,
			"show": false
		}]
	}`)
	var graph Graph
	err := json.Unmarshal(data, &graph)
	if err != nil {
		t.Fatalf("Graph.UnmarshalJSON returned error %s", err)
	}

	expected := NewGraph()
	leftMin := field.ForceString("0")
	expected.YAxes = GraphYaxesOptions{
		Left: GraphYAxis{
			Format:  "dtdurations",
			Label:   "label1",
			LogBase: 1,
			Max:     nil,
			Min:     &leftMin,
			Show:    true,
		},
		Right: GraphYAxis{
			Format:  "short",
			Label:   "",
			LogBase: 10,
			Max:     nil,
			Min:     nil,
			Show:    false,
		},
	}
	if !reflect.DeepEqual(expected, &graph) {
		t.Errorf("Graph.UnmarshalJSON: %s", pretty.Diff(expected, &graph))
	}
}
