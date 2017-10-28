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

	"github.com/spoof/go-grafana/pkg/field"
)

// Graph represents Graph panel
type Graph struct {
	YAxes GraphYaxesOptions `json:"yaxes"`

	generalOptions GeneralOptions
	//queriesOptions QueriesOptions
	queries []Query
}

// NewGraph creates new Graph panel.
func NewGraph() *Graph {
	return &Graph{}
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

func (axe *GraphYaxesOptions) MarshalJSON() ([]byte, error) {
	axes := []GraphYAxis{axe.Left, axe.Right}
	return json.Marshal(axes)
}
func (axe *GraphYaxesOptions) UnmarshalJSON(data []byte) error {
	var axes []GraphYAxis
	if err := json.Unmarshal(data, &axes); err != nil {
		return err
	}
	if len(axes) < 2 {
		return errors.New("Axes should be 2")
	}

	axe.Left = axes[0]
	axe.Right = axes[1]

	return nil
}

type GraphYAxis struct {
	Format  string             `json:"format"`
	Label   string             `json:"label,omitempty"`
	LogBase int                `json:"logBase"`
	Max     *field.ForceString `json:"max"`
	Min     *field.ForceString `json:"min"`
	Show    bool               `json:"show"`
}
