package panel

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/guregu/null"
	"github.com/kr/pretty"
	"github.com/spoof/go-grafana/pkg/field"
	jsontools "github.com/spoof/go-grafana/pkg/json"
)

func TestGraph_MarshalJSON(t *testing.T) {
	p := NewGraph()
	leftMin := field.ForceString("0")
	rightMax := field.ForceString("10")

	p.XAxis.Buckets = null.IntFrom(1)
	p.XAxis.Mode = "series"
	p.XAxis.Show = true
	p.XAxis.Values = []string{"max"}
	p.YAxes = GraphYaxesOptions{
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
	p.Legend.AsTable = true
	p.Legend.Avg = true
	p.Legend.Current = true
	p.Legend.HideEmpty = true
	p.Legend.HideZero = true
	p.Legend.Max = true
	p.Legend.Min = true
	p.Legend.AtRight = true
	p.Legend.Show = true
	p.Legend.Width = null.IntFromPtr(nil)
	p.Legend.Total = true
	p.Legend.Values = true
	got, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		t.Fatalf("Graph.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"decimals": null,
		"xaxis": {
			"buckets": 1,
			"mode": "series",
			"show": true,
			"values": ["max"]
		},
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
		}],
		"legend": {
			"alignAsTable": true,
			"avg": true,
			"current": true,
			"hideEmpty": true,
			"hideZero": true,
			"max": true,
			"min": true,
			"rightSide": true,
			"show": true,
			"sideWidth": null,
			"total": true,
			"values": true
		}
	}`)
	if eq, err := jsontools.BytesEqual(expected, got); err != nil {
		t.Fatalf("Graph.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("Graph.MarshalJSON:\ngot %s\nwant: %s", expected, got)
	}
}
func TestGraph_UnmarshalJSON(t *testing.T) {
	data := []byte(`{
		"decimals": 3,
		"xaxis": {
			"buckets": null,
			"mode": "histogram",
			"show": true,
			"value": []
		},
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
		}],
		"legend": {
			"alignAsTable": true,
			"avg": true,
			"current": true,
			"hideEmpty": true,
			"hideZero": true,
			"max": true,
			"min": true,
			"rightSide": true,
			"show": true,
			"sideWidth": 100,
			"total": true,
			"values": true
		}
	}`)
	var graph Graph
	err := json.Unmarshal(data, &graph)
	if err != nil {
		t.Fatalf("Graph.UnmarshalJSON returned error %s", err)
	}

	expected := NewGraph()
	expected.XAxis.Buckets = null.IntFromPtr(nil)
	expected.XAxis.Mode = "histogram"
	expected.XAxis.Show = true
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
	expected.Legend.AsTable = true
	expected.Legend.Avg = true
	expected.Legend.Current = true
	expected.Legend.HideEmpty = true
	expected.Legend.HideZero = true
	expected.Legend.Max = true
	expected.Legend.Min = true
	expected.Legend.AtRight = true
	expected.Legend.Show = true
	expected.Legend.Width = null.IntFrom(100)
	expected.Legend.Total = true
	expected.Legend.Values = true
	expected.Decimals = null.IntFrom(3)
	if !reflect.DeepEqual(expected, &graph) {
		t.Errorf("Graph.UnmarshalJSON: %s", pretty.Diff(expected, &graph))
	}
}
