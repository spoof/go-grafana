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
	// Legend
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
	// Draw options
	p.DrawOptions.Bars = true
	p.DrawOptions.Lines = true
	p.DrawOptions.Points = true
	p.DrawOptions.Fill = 2
	p.DrawOptions.LineWidth = 3
	p.DrawOptions.PointRadius = 4
	p.DrawOptions.Staircase = true
	p.Tooltip.Shared = true
	p.Tooltip.Sort = Desc
	p.Tooltip.StackedValue = Individual
	p.DrawOptions.Stack = true
	p.DrawOptions.NullValue = NullAsZeroPointMode
	p.SeriesOverrides = []GraphSeriesOverride{
		GraphSeriesOverride{
			Alias:         "/yo/",
			Bars:          boolRef(true),
			Color:         "#E5A8E2",
			Dashes:        boolRef(true),
			DashLength:    uintRef(10),
			DashSpace:     uintRef(2),
			FillBelowTo:   "x",
			Legend:        boolRef(true),
			Lines:         boolRef(true),
			LineFill:      uintRef(2),
			LineWidth:     uintRef(2),
			NullPointMode: NullAsZeroPointMode,
			Points:        boolRef(true),
			PointRadius:   uintRef(2),
			Staircase:     boolRef(true),
			Stack:         "A",
			YAxis:         uintRef(1),
			ZIndex:        intRef(-2),
			Transform:     "negative-y",
		},
	}
	p.Thresholds = []Threshold{
		Threshold{
			Mode:  CustomThresholdMode,
			Fill:  true,
			Line:  true,
			Op:    GreaterOp,
			Value: 1,
			Color: "rgb(77, 58, 58)",
		},
	}

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
		},
		"stack": true,
		"nullPointMode": "null as zero",
		"bars": true,
		"lines": true,
		"points": true,
		"fill": 2,
		"linewidth": 3,
		"steppedLine": true,
		"pointradius": 4,
		"tooltip": {
			"shared": true,
			"sort": 2,
			"value_type": "individual"
		},
		"seriesOverrides": [{
			"alias": "/yo/",
			"bars": true,
			"color": "#E5A8E2",
			"dashes": true,
			"dashLength": 10,
			"spaceLength": 2,
			"fillBelowTo": "x",
			"legend": true,
			"lines": true,
			"fill": 2,
			"linewidth": 2,
			"nullPointMode": "null as zero",
			"points": true,
			"pointradius": 2,
			"steppedLine": true,
			"stack": "A",
			"yaxis": 1,
			"zindex": -2,
			"transform": "negative-y"
		}],
		"thresholds": [{
			"colorMode": "custom",
			"fill": true,
			"line": true,
			"lineColor": "rgb(77, 58, 58)",
			"op": "gt",
			"value": 1
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
		},
		"stack": true,
		"nullPointMode": "connected",
		"bars": true,
		"lines": true,
		"points": true,
		"fill": 2,
		"linewidth": 3,
		"steppedLine": true,
		"pointradius": 4,
		"tooltip": {
			"shared": true,
			"sort": 1,
			"value_type": "cumulative"
		},
		"seriesOverrides": [{
			"alias": "/yo/",
			"dashes": true,
			"bars": true,
			"color": "#E5A8E2",
			"dashes": true,
			"dashLength": 10,
			"fillBelowTo": "x",
			"legend": true,
			"lines": true,
			"fill": 2,
			"linewidth": 2,
			"nullPointMode": "null as zero",
			"points": true,
			"pointradius": 2,
			"steppedLine": true,
			"spaceLength": 2,
			"stack": "A",
			"yaxis": 1,
			"zindex": -2,
			"transform": "negative-y"
		}],
		"thresholds": [{
			"colorMode": "custom",
			"fill": true,
			"line": true,
			"lineColor": "rgb(77, 58, 58)",
			"op": "gt",
			"value": 1
		}]
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

	// Draw options
	expected.DrawOptions.Bars = true
	expected.DrawOptions.Lines = true
	expected.DrawOptions.Points = true
	expected.DrawOptions.Fill = 2
	expected.DrawOptions.LineWidth = 3
	expected.DrawOptions.PointRadius = 4
	expected.DrawOptions.Staircase = true
	expected.Tooltip.Shared = true
	expected.Tooltip.Sort = Asc
	expected.Tooltip.StackedValue = Cumulative
	expected.DrawOptions.Stack = true
	expected.DrawOptions.NullValue = ConnectedNullPointMode

	// Series Overrides
	expected.SeriesOverrides = []GraphSeriesOverride{
		GraphSeriesOverride{
			Alias:         "/yo/",
			Bars:          boolRef(true),
			Color:         "#E5A8E2",
			Dashes:        boolRef(true),
			DashLength:    uintRef(10),
			DashSpace:     uintRef(2),
			FillBelowTo:   "x",
			Legend:        boolRef(true),
			Lines:         boolRef(true),
			LineFill:      uintRef(2),
			LineWidth:     uintRef(2),
			NullPointMode: NullAsZeroPointMode,
			Points:        boolRef(true),
			PointRadius:   uintRef(2),
			Staircase:     boolRef(true),
			Stack:         "A",
			YAxis:         uintRef(1),
			ZIndex:        intRef(-2),
			Transform:     "negative-y",
		},
	}
	// Thresholds
	expected.Thresholds = []Threshold{
		Threshold{
			Mode:  CustomThresholdMode,
			Fill:  true,
			Line:  true,
			Op:    GreaterOp,
			Value: 1,
			Color: "rgb(77, 58, 58)",
		},
	}

	if !reflect.DeepEqual(expected, &graph) {
		t.Errorf("Graph.UnmarshalJSON: %s", pretty.Diff(expected, &graph))
	}
}

func boolRef(b bool) *bool {
	return &b
}

func intRef(i int) *int {
	return &i
}

func uintRef(i uint) *uint {
	return &i
}
