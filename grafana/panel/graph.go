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

package panel

import (
	"encoding/json"
	"errors"

	"github.com/guregu/null"
	"github.com/spoof/go-grafana/pkg/field"
)

type graphXAxisMode string

const (
	graphXAxisHistogram graphXAxisMode = "histogram"
	graphXAxisSeries    graphXAxisMode = "series"
	graphXAxisTime      graphXAxisMode = "time"
)

type sortType uint

const (
	NoSort sortType = iota
	Asc
	Desc
)

type stackedValueType string

const (
	Individual stackedValueType = "individual"
	Cumulative stackedValueType = "cumulative"
)

type nullPointMode string

const (
	ConnectedNullPointMode nullPointMode = "connected"
	NullNullPointMode      nullPointMode = "null"
	NullAsZeroPointMode    nullPointMode = "null as zero"
)

// Graph represents Graph panel
type Graph struct {
	// Axes
	YAxes GraphYaxesOptions `json:"yaxes"`
	XAxis struct {
		Buckets null.Int       `json:"buckets,omitempty"`
		Mode    graphXAxisMode `json:"mode"`           // TODO: add validaition here: histogram/series/time
		Name    *string        `json:"name,omitempty"` // it's seems that it's not used anymore
		Show    bool           `json:"show"`
		Values  []string       `json:"values"` // TODO: actually it's only single value here. Need custom type
	} `json:"xaxis"`

	// Legend
	Legend struct {
		AsTable   bool     `json:"alignAsTable"`
		Avg       bool     `json:"avg"`
		Current   bool     `json:"current"`
		HideEmpty bool     `json:"hideEmpty"`
		HideZero  bool     `json:"hideZero"`
		Max       bool     `json:"max"`
		Min       bool     `json:"min"`
		AtRight   bool     `json:"rightSide"`
		Show      bool     `json:"show"`
		Width     null.Int `json:"sideWidth"`
		Total     bool     `json:"total"`
		Values    bool     `json:"values"`
	} `json:"legend"`
	Decimals null.Int `json:"decimals"`

	// Display
	DrawOptions
	Tooltip struct {
		// All series or single
		Shared       bool             `json:"shared"`
		Sort         sortType         `json:"sort"`
		StackedValue stackedValueType `json:"value_type"`
	} `json:"tooltip"`

	SeriesOverrides []GraphSeriesOverride `json:"seriesOverrides"`

	generalOptions GeneralOptions
	//queriesOptions QueriesOptions
	queries []Query
}

// NewGraph creates new Graph panel.
func NewGraph() *Graph {
	return &Graph{}
}

type DrawOptions struct {
	Bars   bool `json:"bars"`
	Lines  bool `json:"lines"`
	Points bool `json:"points"`

	// Options
	Fill        uint `json:"fill"`
	LineWidth   uint `json:"linewidth"`
	Staircase   bool `json:"steppedLine"`
	PointRadius uint `json:"pointradius,omitempty"`

	Stack     bool          `json:"stack"`
	NullValue nullPointMode `json:"nullPointMode"`
}

// GeneralOptions implements grafana.Panel interface
func (p *Graph) GeneralOptions() *GeneralOptions {
	return &p.generalOptions
}

// Queries implements Queryable interface
func (p *Graph) Queries() *[]Query {
	return &p.queries
}

type GraphYaxesOptions struct {
	Left  GraphYAxis
	Right GraphYAxis
}

func (y *GraphYaxesOptions) MarshalJSON() ([]byte, error) {
	axes := []GraphYAxis{y.Left, y.Right}
	return json.Marshal(axes)
}

func (y *GraphYaxesOptions) UnmarshalJSON(data []byte) error {
	var axes []GraphYAxis
	if err := json.Unmarshal(data, &axes); err != nil {
		return err
	}
	if len(axes) < 2 {
		return errors.New("Axes should be 2")
	}

	y.Left = axes[0]
	y.Right = axes[1]

	return nil
}

type GraphYAxis struct {
	Format  string             `json:"format"` // TODO: replace with custom type with default value "short".
	Label   string             `json:"label,omitempty"`
	LogBase int                `json:"logBase"` // TODO: default value should be 1 (linear)
	Max     *field.ForceString `json:"max"`
	Min     *field.ForceString `json:"min"`
	Show    bool               `json:"show"`
}

type GraphSeriesOverride struct {
	Alias         string        `json:"alias"`
	Bars          *bool         `json:"bars"`
	Color         string        `json:"color,omitempty"`
	Dashes        *bool         `json:"dashes"`
	DashLength    *uint         `json:"dashLength"`
	DashSpace     *uint         `json:"spaceLength"`
	FillBelowTo   string        `json:"fillBelowTo,omitempty"`
	Legend        *bool         `json:"legend"`
	Lines         *bool         `json:"lines"`
	LineFill      *uint         `json:"fill"`
	LineWidth     *uint         `json:"linewidth"`
	NullPointMode nullPointMode `json:"nullPointMode,omitempty"`
	Points        *bool         `json:"points,omitempty"`
	PointRadius   *uint         `json:"pointradius"`
	Staircase     *bool         `json:"steppedLine"`
	Stack         string        `json:"stack"` // could be bool (true/false)
	YAxis         *uint         `json:"yaxis"`
	ZIndex        *int          `json:"zindex"`
	Transform     string        `json:"transform"` // always 'negative-y'
}
