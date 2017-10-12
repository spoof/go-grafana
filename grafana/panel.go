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
	Datasource string `json:"datasource,omitempty"`
	Queries    []Query
}

func (o *QueriesOptions) UnmarshalJSON(data []byte) error {
	type JSONOptions QueriesOptions

	var queries []*probeQuery
	jo := struct {
		*JSONOptions
		Queries *[]*probeQuery `json:"targets,omitempty"`
	}{
		JSONOptions: (*JSONOptions)(o),
		Queries:     &queries,
	}
	if err := json.Unmarshal(data, &jo); err != nil {
		return err
	}

	o.Queries = []Query{}
	for _, q := range queries {
		// TODO: queies shouldn't be nil in future. This check will be obsolete
		if q.query == nil {
			continue
		}
		o.Queries = append(o.Queries, q.query)
	}

	return nil
}

type probeQuery struct {
	// PrometheusQuery
	IntervalFactor *uint   `json:"intervalFactor"`
	Expression     *string `json:"expr"`

	// GraphiteQuery fields
	Target *string `json:"target"`

	query Query
}

func (q *probeQuery) UnmarshalJSON(data []byte) error {
	type JSONQuery probeQuery
	var jq JSONQuery
	if err := json.Unmarshal(data, &jq); err != nil {
		return err
	}

	var query Query
	if jq.Expression != nil && jq.IntervalFactor != nil {
		query = new(PrometheusQuery)
	} else if jq.Target != nil {
		query = new(GraphiteQuery)
	}

	// TODO: Initialize Unknown query here instead
	if query == nil {
		return nil
	}

	if err := json.Unmarshal(data, &query); err != nil {
		return err
	}

	q.query = query
	return nil
}

// commonQuery is a single query options that common to all datasources.
type commonQuery struct {
	Datasource *string `json:"datasource"`
	RefID      string  `json:"refid"`
}

// PrometheusQuery is query specific options for Prometheus datasource.
type PrometheusQuery struct {
	commonQuery

	IntervalFactor uint `json:"intervalFactor"`
	// Interval       uint   `json:"interval"`
	// FIXME: Interval can be a string. We need to convert it to int
	Interval     string `json:"interval"`
	Format       string `json:"format"`
	Expression   string `json:"expr"`
	LegendFormat string `json:"legendFormat"`
	Step         uint   `json:"step"`
}

func (q *PrometheusQuery) RefID() string {
	return q.commonQuery.RefID
}

func (q *PrometheusQuery) Datasource() string {
	return *q.commonQuery.Datasource
}

// GraphiteQuery is query specific options for Graphite datasource.
type GraphiteQuery struct {
	commonQuery

	Target string `json:"target"`
}

func (q *GraphiteQuery) RefID() string {
	return q.commonQuery.RefID
}

func (q *GraphiteQuery) Datasource() string {
	return *q.commonQuery.Datasource
}

type Query interface {
	RefID() string
	Datasource() string
}
