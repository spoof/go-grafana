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
	"time"

	"github.com/spoof/go-grafana/grafana/panel"
	panelQuery "github.com/spoof/go-grafana/grafana/query"
	"github.com/spoof/go-grafana/pkg/field"
)

type (
	// DashboardID is an ID type of Dashboard
	DashboardID    uint64
	dashboardStyle string
)

const (
	dashboardDarkStyle  dashboardStyle = "dark"
	dashboardLightStyle dashboardStyle = "light"
)

type Dashboard struct {
	ID            DashboardID `json:"-"`
	Version       uint64      `json:"-"`
	SchemaVersion int         `json:"schemaVersion"`

	Editable     bool           `json:"editable"`
	GraphTooltip uint8          `json:"graphTooltip"`
	HideControls bool           `json:"hideControls"`
	Rows         []*Row         `json:"rows"`
	Style        dashboardStyle `json:"style"`
	Timezone     string         `json:"timezone"`
	Title        string         `json:"title"`
	Tags         *field.Tags    `json:"tags"`

	Meta *DashboardMeta `json:"-"`
}

// NewDashboard creates new Dashboard.
func NewDashboard(title string) *Dashboard {
	return &Dashboard{
		Title:         title,
		Editable:      true,
		SchemaVersion: 14,
		Style:         dashboardDarkStyle,
		Tags:          field.NewTags(),
	}
}

// MarshalJSON implements json.Marshaler interface
func (d *Dashboard) MarshalJSON() ([]byte, error) {
	panelID := 1
	rows := make([]*Row, len(d.Rows))
	for i, r := range d.Rows {
		rr := *r

		panels := make([]Panel, len(r.Panels))
		for j, p := range r.Panels {
			panels[j] = &probePanel{
				ID:    uint(panelID),
				panel: p,
			}
			panelID++
		}
		rr.Panels = panels
		rows[i] = &rr
	}

	type JSONDashboard Dashboard
	jd := &struct {
		JSONDashboard
		Rows []*Row   `json:"rows"`
		Tags []string `json:"tags"`
	}{
		JSONDashboard: (JSONDashboard)(*d),
		Rows:          rows,
		Tags:          d.Tags.Value(),
	}
	return json.Marshal(jd)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Dashboard) UnmarshalJSON(data []byte) error {
	type JSONDashboard Dashboard
	inDashboard := struct {
		*JSONDashboard
		ID      *DashboardID `json:"id"`
		Version *uint64      `json:"version"`

		Tags []string       `json:"tags"`
		Meta *DashboardMeta `json:"meta"`
	}{
		JSONDashboard: (*JSONDashboard)(d),
		ID:            &d.ID,
		Version:       &d.Version,
		Meta:          d.Meta,
	}
	if err := json.Unmarshal(data, &inDashboard); err != nil {
		return err
	}

	d.Tags = field.NewTags(inDashboard.Tags...)

	return nil
}

type DashboardMeta struct {
	Slug    string `json:"slug"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	CanEdit bool `json:"canEdit"`
	CanSave bool `json:"canSave"`
	CanStar bool `json:"canStar"`

	Created   time.Time `json:"created"`
	CreatedBy string    `json:"createdBy"`
	Expires   time.Time `json:"expires"`
	Updated   time.Time `json:"updated"`
	UpdatedBy string    `json:"updatedBy"`
}

func (dm *DashboardMeta) String() string {
	return Stringify(dm)
}

// Row is panel's row
type Row struct {
	Collapsed bool              `json:"collapse"`
	Editable  bool              `json:"editable"`
	Height    field.ForceString `json:"height"`
	Panels    []Panel           `json:"panels"`
	RepeatFor string            `json:"repeat"` // repeat row for given variable
	ShowTitle bool              `json:"showTitle"`
	Title     string            `json:"title"`
	TitleSize string            `json:"titleSize"` // TODO: validation: h1-h6
}

// NewRow creates new Row with somw defaults.
func NewRow() *Row {
	return &Row{Editable: true}
}

// UnmarshalJSON implements json.Unmarshaler interface
func (r *Row) UnmarshalJSON(data []byte) error {
	type JSONRow Row
	jr := struct {
		*JSONRow
		Panels []probePanel `json:"panels"`
	}{
		JSONRow: (*JSONRow)(r),
	}

	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}

	panels := make([]Panel, len(jr.Panels))
	for i, p := range jr.Panels {
		panels[i] = p.panel
	}
	r.Panels = panels
	return nil
}

//Panel represents Dashboard's panel
type Panel interface {
	GeneralOptions() *panel.GeneralOptions
}

type panelType string

const (
	textPanelType       panelType = "text"
	singlestatPanelType panelType = "singlestat"
	graphPanelType      panelType = "graph"
)

type probePanel struct {
	ID   uint      `json:"id"`
	Type panelType `json:"type"`

	panel Panel
}

func (p *probePanel) GeneralOptions() *panel.GeneralOptions {
	return p.panel.GeneralOptions()
}

func (p *probePanel) UnmarshalJSON(data []byte) error {
	type JSONPanel probePanel
	jp := struct {
		*JSONPanel
	}{
		JSONPanel: (*JSONPanel)(p),
	}
	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}

	var pp Panel
	switch jp.Type {
	case textPanelType:
		pp = new(panel.Text)
	case singlestatPanelType:
		pp = new(panel.Singlestat)
	case graphPanelType:
		pp = new(panel.Graph)
	default:
		return nil
	}

	if err := json.Unmarshal(data, pp); err != nil {
		return err
	}

	// Unmarshal general options
	var generalOptions panel.GeneralOptions
	if err := json.Unmarshal(data, &generalOptions); err != nil {
		return err
	}
	gOpts := pp.GeneralOptions()
	*gOpts = generalOptions

	// Unmarshal queries
	var queriesOpts queriesOptions
	if err := json.Unmarshal(data, &queriesOpts); err != nil {
		return err
	}
	if queryablePanel, ok := pp.(QueryablePanel); ok {
		queriesPtr := queryablePanel.Queries()
		newQueries := []panel.Query{}
		for _, q := range queriesOpts.Queries {
			if q.query == nil {
				continue
			}
			newQueries = append(newQueries, q.query)
		}
		*queriesPtr = newQueries
	}

	p.panel = pp
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (p *probePanel) MarshalJSON() ([]byte, error) {
	type JSONPanel probePanel
	jp := struct {
		*JSONPanel

		*panel.Text
		*panel.Singlestat
		*panel.Graph

		*panel.GeneralOptions
		*queriesOptions
	}{
		JSONPanel:      (*JSONPanel)(p),
		GeneralOptions: p.GeneralOptions(),
	}

	switch v := p.panel.(type) {
	case *panel.Text:
		jp.Text = v
		jp.Type = textPanelType
	case *panel.Singlestat:
		jp.Singlestat = v
		jp.Type = singlestatPanelType
	case *panel.Graph:
		jp.Graph = v
		jp.Type = graphPanelType
	}

	if qp, ok := p.panel.(QueryablePanel); ok {
		// Determine do each query uses its own datassource or not
		isOwnDatasource := false
		queries := *qp.Queries()
		for i := 0; i < len(queries)-1; i++ {
			if queries[i].Datasource() != queries[i+1].Datasource() {
				isOwnDatasource = true
			}
		}

		probeQueries := make([]probeQuery, len(queries))
		for i, q := range queries {
			pq := probeQuery{
				RefID: makeRefID(i),
				query: q,
			}

			if isOwnDatasource {
				pq.Datasource = q.Datasource()
			}
			probeQueries[i] = pq
		}

		var datasource string
		if isOwnDatasource {
			datasource = mixedDatasource
		} else {
			if len(queries) > 0 {
				datasource = queries[0].Datasource()
			}
		}

		jp.queriesOptions = &queriesOptions{
			Queries:    probeQueries,
			Datasource: datasource,
		}
	}

	return json.Marshal(jp)
}

// QueryablePanel is interface for panels that supports quering metrics from datasources.
type QueryablePanel interface {
	Queries() *[]panel.Query
}

const mixedDatasource = "-- Mixed --"

type queriesOptions struct {
	Datasource string       `json:"datasource,omitempty"`
	Queries    []probeQuery `json:"targets"`
}

// probeQuery is an auxiliary entity thats purpose to manage marshaling and unmarshal of panel's query into concrete
// types.
type probeQuery struct {
	RefID      string `json:"refid"`
	Datasource string `json:"datasource,omitempty"`

	query panel.Query
}

// UnmarshalJSON implements json.Unmarshaler interface
func (q *probeQuery) UnmarshalJSON(data []byte) error {
	type JSONQuery probeQuery
	jq := struct {
		*JSONQuery

		// Prometheus query fields
		IntervalFactor *uint   `json:"intervalFactor"`
		Expression     *string `json:"expr"`

		// Graphite queryfields
		Target *string `json:"target"`
	}{
		JSONQuery: (*JSONQuery)(q),
	}
	if err := json.Unmarshal(data, &jq); err != nil {
		return err
	}

	// There is no any information about query type in Grafana's JSON object. Further more, some queries uses the same
	// fields. Thus we need to use some heurisitcs to map json fields into our query types properly.
	// This heurisitcs based on searching specific for query type fields in JSON data.
	var query panel.Query
	if jq.Expression != nil && jq.IntervalFactor != nil {
		query = new(panelQuery.Prometheus)
	} else if jq.Target != nil {
		query = new(panelQuery.Graphite)
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

// MarshalJSON implements json.Marshaler interface
func (q *probeQuery) MarshalJSON() ([]byte, error) {
	type JSONQuery probeQuery
	jq := struct {
		*JSONQuery
		*panelQuery.Prometheus
		*panelQuery.Graphite
	}{
		JSONQuery: (*JSONQuery)(q),
	}

	switch v := q.query.(type) {
	case *panelQuery.Prometheus:
		jq.Prometheus = v
	case *panelQuery.Graphite:
		jq.Graphite = v
	}

	return json.Marshal(jq)
}

// makeRefID returns symbolic ID for given index.
// TODO: It has very rough implementation. Needs refactoring.
func makeRefID(index int) string {
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	var id string
	if index >= len(letters) {
		id += makeRefID(index % len(letters))
	} else {
		id = string(letters[index])
	}

	var result string
	for _, v := range id {
		result = string(v) + result
	}
	return result
}
