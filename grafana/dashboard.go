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
	"github.com/spoof/go-grafana/pkg/field"
)

type (
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
	tags         *field.Tags

	Meta *DashboardMeta `json:"-"`
}

// NewDashboard creates new Dashboard.
func NewDashboard(title string) *Dashboard {
	return &Dashboard{
		Title:         title,
		Editable:      true,
		SchemaVersion: 14,
		Style:         dashboardDarkStyle,
		tags:          field.NewTags(),
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
		Tags:          d.tags.Value(),
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

	d.tags = field.NewTags(inDashboard.Tags...)

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

type Panel interface {
	GeneralOptions() *panel.GeneralOptions
}

type panelType string

const (
	textPanel       panelType = "text"
	singlestatPanel panelType = "singlestat"
	graphPanel      panelType = "graph"
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
	case textPanel:
		pp = new(panel.TextPanel)
	case singlestatPanel:
		pp = new(panel.SinglestatPanel)
	case graphPanel:
		pp = new(panel.GraphPanel)
	default:
		return nil
	}

	if err := json.Unmarshal(data, pp); err != nil {
		return err
	}

	var generalOptions panel.GeneralOptions
	if err := json.Unmarshal(data, &generalOptions); err != nil {
		return err
	}
	gOpts := pp.GeneralOptions()
	*gOpts = generalOptions

	//var queriesOptions PanelQueriesOptions
	//if err := json.Unmarshal(data, &queriesOptions); err != nil {
	//	return err
	//}

	p.panel = pp
	return nil
}

// MarshalJSON implements json.Marshaler interface
func (p *probePanel) MarshalJSON() ([]byte, error) {
	type JSONPanel probePanel
	jp := struct {
		*JSONPanel

		*panel.TextPanel
		*panel.SinglestatPanel
		*panel.GraphPanel

		*panel.GeneralOptions
	}{
		JSONPanel:      (*JSONPanel)(p),
		GeneralOptions: p.GeneralOptions(),
	}

	switch v := p.panel.(type) {
	case *panel.TextPanel:
		jp.TextPanel = v
		jp.Type = textPanel
	case *panel.SinglestatPanel:
		jp.SinglestatPanel = v
		jp.Type = singlestatPanel
	case *panel.GraphPanel:
		jp.GraphPanel = v
		jp.Type = graphPanel
	}
	return json.Marshal(jp)
}

type panelGeneralOptions struct {
	id        uint
	panelType panelType
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

	//var query Query
	//if jq.Expression != nil && jq.IntervalFactor != nil {
	//	query = new(PrometheusQuery)
	//} else if jq.Target != nil {
	//	query = new(GraphiteQuery)
	//}

	// TODO: Initialize Unknown query here instead
	//if query == nil {
	//	return nil
	//}

	//if err := json.Unmarshal(data, &query); err != nil {
	//	return err
	//}

	//q.query = query
	return nil
}

type Query interface {
	RefID() string
	Datasource() string
	//commonOptions() *query.commonQuery
}

// PanelQueriesOptions is a part of panel that placed in 'Metrics' tab. It represents set of panel queries.
type PanelQueriesOptions struct {
	Datasource string  `json:"datasource,omitempty"`
	Queries    []Query `json:"targets"`
}

func (o *PanelQueriesOptions) UnmarshalJSON(data []byte) error {
	type JSONOptions PanelQueriesOptions

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

	//o.Queries = []Query{}
	//for _, q := range queries {
	//	// TODO: queies shouldn't be nil in future. This check will be obsolete
	//	if q.query == nil {
	//		continue
	//	}
	//	o.Queries = append(o.Queries, q.query)
	//}

	return nil
}

// MarshalJSON implements encoding/json.Marshaler
func (o *PanelQueriesOptions) MarshalJSON() ([]byte, error) {
	type JSONOptions PanelQueriesOptions
	jo := (*JSONOptions)(o)

	// FIXME: add checking for uniqueness of refids
	//for i, q := range jo.Queries {
	//	if q.commonOptions().RefID != "" {
	//		continue
	//	}
	//	q.commonOptions().RefID = makeRefID(i)
	//}

	// TODO: if there are a several types of datasources we need to set 'main' datasouce to "Mixed"

	return json.Marshal(jo)
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
