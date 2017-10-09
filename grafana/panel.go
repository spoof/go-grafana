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
)

type Panel interface {
	GeneralOptions() *PanelGeneralOptions
	QueriesOptions() *QueriesOptions
}

type panelType string

const (
	textPanel       panelType = "text"
	singlestatPanel panelType = "singlestat"
	graphPanel      panelType = "graph"
)

type probePanel struct {
	Type panelType `json:"type"`
	Panel
}

func (p *probePanel) UnmarshalJSON(data []byte) error {
	type JSONPanel probePanel
	var jp JSONPanel
	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}

	var panel Panel
	switch jp.Type {
	case textPanel:
		panel = new(TextPanel)
	case singlestatPanel:
		panel = new(SinglestatPanel)
	case graphPanel:
		panel = new(GraphPanel)
	default:
		return nil
	}

	if err := json.Unmarshal(data, panel); err != nil {
		return err
	}

	p.Panel = panel
	return nil
}

type PanelGeneralOptions struct {
	id        uint
	panelType panelType

	Description string       `json:"description"`
	Height      forceString  `json:"height"`
	Links       []*PanelLink `json:"links"`
	MinSpan     uint         `json:"minSpan"` // TODO: valid values: 1-12
	Span        uint         `json:"span"`    // TODO: valid values: 1-12
	Title       string       `json:"title"`
	Transparent bool         `json:"transparent"`
}

type panelLinkType string

const (
	PanelLinkAbsolute  panelLinkType = "absolute"
	PanelLinkDashboard panelLinkType = "dashboard"
)

type PanelLink struct {
	IncludeVars  bool          `json:"includeVars"`
	KeepTime     bool          `json:"keepTime"`
	Params       string        `json:"params"`
	OpenInNewTab bool          `json:"targetBlank"`
	Type         panelLinkType `json:"type"`

	// type=absolute
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`

	// type=dashboard
	DashboardURI string `json:"dashUri,omitempty"`   // TODO: validation. should be valid dashboard
	Dashboard    string `json:"dashboard,omitempty"` // actually it's title. Autofilled from dashboard title
}

// NewPanelLink creates new PanelLink
func NewPanelLink(linkType panelLinkType) *PanelLink {
	return &PanelLink{
		Type: linkType,
	}
}

// QueriesOptions is a part of panel that placed in 'Metrics' tab. It represents set of panel queries.
type QueriesOptions struct {
	Datasource string   `json:"datasource,omitempty"`
	Queries    []*Query `json:"targets,omitempty"`
}

// Query is a single query options that common to all datasources.
type Query struct {
	Datasource *string `json:"datasource"`
	RefID      string  `json:"refid"`
}

// PrometheusQuery is query specific options for Prometheus datasource.
type PrometheusQuery struct {
	Query

	IntervalFactor uint   `json:"intervalFactor"`
	Interval       uint   `json:"interval"`
	Format         string `json:"format"`
	Expression     string `json:"expr"`
	LegendFormat   string `json:"legendFormat"`
	Step           uint   `json:"step"`
}

// GraphiteQuery is query specific options for Graphite datasource.
type GraphiteQuery struct {
	Query

	Target string `json:"target"`
}
