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

// Graphite is query specific options for Graphite datasource.
type Graphite struct {
	Target     string `json:"target"`
	TargetFull string `json:"targetFull,omitempty"`

	datasource string
}

// NewGraphite creates new instance of Graphite query
func NewGraphite(datasourceName string) *Graphite {
	return &Graphite{
		datasource: datasourceName,
	}
}

// Datasource implements panel.Query interface
func (q *Graphite) Datasource() string {
	return q.datasource
}
