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

package panel_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/guregu/null"
	"github.com/kr/pretty"
	"github.com/spoof/go-grafana/grafana/panel"
	jsontools "github.com/spoof/go-grafana/pkg/json"
)

func TestSinglestat_MarshalJSON(t *testing.T) {
	p := panel.NewSinglestat()
	p.ValueName = "max"
	p.ValueFontSize = "50%"
	p.Prefix = "prefix"
	p.PrefixFontSize = "75%"
	p.Postfix = "postfix"
	p.PostfixFontSize = "25%"
	p.Format = "currencyJPY"
	p.ColorBackground = true
	p.ColorValue = true
	p.Thresholds = "20,80"
	p.Colors = []string{"rgba(10, 172, 45, 0.97)", "rgba(20, 172, 45, 0.97)", "rgba(30, 172, 45, 0.97)"}

	p.SparkLine.Show = true
	p.SparkLine.FullHeight = true
	p.SparkLine.LineColor = "rgb(31, 120, 193)"
	p.SparkLine.FillColor = "rgba(50, 61, 71, 0.18)"

	p.Gauge.Show = true
	p.Gauge.MaxValue = 100
	p.Gauge.MinValue = 50
	p.Gauge.ThresholdLabels = true
	p.Gauge.ThresholdMarkers = true

	options := p.GeneralOptions()
	options.Title = "New Panel"
	options.Description = "Panel Description"
	options.Height = "250px"
	options.MinSpan = 1
	options.Span = 12
	options.Transparent = true

	// Time range
	p.TimeRangeOptions = panel.TimeRangeOptions{
		From:         null.StringFrom("1h"),
		Shift:        null.StringFromPtr(nil),
		HideOverride: true,
	}

	got, err := json.MarshalIndent(p, "", "\t\t")
	if err != nil {
		t.Fatalf("Singlestat.MarshalJSON returned error %s", err)
	}
	expected := []byte(`{
		"valueName": "max",
		"valueFontSize": "50%",
		"prefix": "prefix",
		"prefixFontSize": "75%",
		"postfix": "postfix",
		"postfixFontSize": "25%",
		"format": "currencyJPY",
		"colorBackground": true,
		"colorValue": true,
		"thresholds": "20,80",
		"colors": [
			"rgba(10, 172, 45, 0.97)",
			"rgba(20, 172, 45, 0.97)",
			"rgba(30, 172, 45, 0.97)"
		],
		"sparkline": {
			"show": true,
			"full": true,
			"lineColor": "rgb(31, 120, 193)",
			"fillColor": "rgba(50, 61, 71, 0.18)"
		},
		"gauge": {
			"show": true,
			"maxValue": 100,
			"minValue": 50,
			"thresholdLabels": true,
			"thresholdMarkers": true
		},

		"timeFrom": "1h",
		"timeShift": null,
		"hideTimeOverride": true
	}`)
	if eq, err := jsontools.BytesEqual(expected, got); err != nil {
		t.Fatalf("Singlestat.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("Singlestat.MarshalJSON: %s", pretty.Diff(expected, &got))
	}
}

func TestSinglestat_UnmarshalJSON(t *testing.T) {
	expected := panel.NewSinglestat()
	expected.ValueName = "max"
	expected.ValueFontSize = "50%"
	expected.Prefix = "prefix"
	expected.PrefixFontSize = "75%"
	expected.Postfix = "postfix"
	expected.PostfixFontSize = "25%"
	expected.Format = "currencyJPY"
	expected.ColorBackground = true
	expected.ColorValue = true
	expected.Thresholds = "20,80"
	expected.Colors = []string{"rgba(10, 172, 45, 0.97)", "rgba(20, 172, 45, 0.97)", "rgba(30, 172, 45, 0.97)"}

	expected.SparkLine.Show = true
	expected.SparkLine.FullHeight = true
	expected.SparkLine.LineColor = "rgb(31, 120, 193)"
	expected.SparkLine.FillColor = "rgba(50, 61, 71, 0.18)"

	expected.Gauge.Show = true
	expected.Gauge.MaxValue = 100
	expected.Gauge.MinValue = 50
	expected.Gauge.ThresholdLabels = true
	expected.Gauge.ThresholdMarkers = true

	// Time range
	expected.TimeRangeOptions = panel.TimeRangeOptions{
		From:         null.StringFrom("1h"),
		Shift:        null.StringFromPtr(nil),
		HideOverride: true,
	}

	data := []byte(`{
		"valueName": "max",
		"valueFontSize": "50%",
		"postfix": "postfix",
		"postfixFontSize": "25%",
		"prefix": "prefix",
		"prefixFontSize": "75%",
		"format": "currencyJPY",
		"colorBackground": true,
		"colorValue": true,
		"thresholds": "20,80",
		"colors": [
		"rgba(10, 172, 45, 0.97)",
		"rgba(20, 172, 45, 0.97)",
		"rgba(30, 172, 45, 0.97)"
		],
		"sparkline": {
			"show": true,
			"full": true,
			"lineColor": "rgb(31, 120, 193)",
			"fillColor": "rgba(50, 61, 71, 0.18)"
		},
		"gauge": {
			"show": true,
			"maxValue": 100,
			"minValue": 50,
			"thresholdLabels": true,
			"thresholdMarkers": true
		},

		"timeFrom": "1h",
		"timeShift": null,
		"hideTimeOverride": true
	}`)
	var got panel.Singlestat
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Singlestat.UnmarshalJSON returned error %s", err)
	}

	if !reflect.DeepEqual(expected, &got) {
		t.Errorf("Singlestat.UnmarshalJSON: %s", pretty.Diff(expected, &got))
	}
}
