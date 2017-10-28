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

package query

// Prometheus is query specific options for Prometheus datasource.
type Prometheus struct {
	IntervalFactor uint `json:"intervalFactor"`
	// Interval       uint   `json:"interval"`
	// FIXME: Interval can be a string. We need to convert it to int
	Interval     string `json:"interval,omitempty"`
	Format       string `json:"format"`
	Expression   string `json:"expr"`
	LegendFormat string `json:"legendFormat,omitempty"`
	Step         uint   `json:"step,omitempty"`

	datasource string
}

// NewPrometheus creates new instance of Prometheus query.
func NewPrometheus(datasourceName string) *Prometheus {
	return &Prometheus{
		datasource: datasourceName,
	}
}

// Datasource implements panel.Query interface
func (q *Prometheus) Datasource() string {
	return q.datasource
}
