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

func TestSinglestatPanel_MarshalJSON(t *testing.T) {
	panel := NewSinglestatPanel()
	panel.ValueName = "max"
	panel.ValueFontSize = "50%"
	panel.Prefix = "prefix"
	panel.PrefixFontSize = "75%"
	panel.Postfix = "postfix"
	panel.PostfixFontSize = "25%"
	panel.Format = "currencyJPY"
	panel.ColorBackground = true
	panel.ColorValue = true
	panel.Thresholds = "20,80"
	panel.Colors = []string{"rgba(10, 172, 45, 0.97)", "rgba(20, 172, 45, 0.97)", "rgba(30, 172, 45, 0.97)"}

	panel.SparkLine.Show = true
	panel.SparkLine.FullHeight = true
	panel.SparkLine.LineColor = "rgb(31, 120, 193)"
	panel.SparkLine.FillColor = "rgba(50, 61, 71, 0.18)"

	panel.Gauge.Show = true
	panel.Gauge.MaxValue = 100
	panel.Gauge.MinValue = 50
	panel.Gauge.ThresholdLabels = true
	panel.Gauge.ThresholdMarkers = true

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
		"description": "Panel Description",
		"height": "250px",
		"links": [],
		"minSpan": 1,
		"span": 12,
		"title": "New Panel",
		"transparent": true,
		"id": 0,
		"type": "singlestat"
	}`)
	if eq, err := JSONBytesEqual(expected, got); err != nil {
		t.Fatalf("SinglestatPanel.MarshalJSON returned error %s", err)
	} else if !eq {
		t.Errorf("SinglestatPanel.MarshalJSON:\nexpected: %s\ngot: %s", expected, got)
	}
}

func TestSinglestatPanel_UnmarshalJSON(t *testing.T) {
	expected := NewSinglestatPanel()
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

	options := expected.GeneralOptions()
	options.Title = "New Panel"
	options.Description = "Panel Description"
	options.Height = "250px"
	options.MinSpan = 1
	options.Span = 12
	options.Transparent = true

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
		"description": "Panel Description",
		"height": "250px",
		"links": [],
		"minSpan": 1,
		"span": 12,
		"title": "New Panel",
		"transparent": true,
		"id": 0,
		"type": "singlestat"
	}`)
	var got SinglestatPanel
	err := json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("SinglestatPanel.UnmarshalJSON returned error %s", err)
	}

	if !reflect.DeepEqual(expected, &got) {
		t.Errorf("SinglestatPanel.MarshalJSON:\nexpected: %#v\ngot: %#v", expected, &got)
	}
}
