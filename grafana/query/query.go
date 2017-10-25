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
	Interval     string `json:"interval,omitempty"`
	Format       string `json:"format"`
	Expression   string `json:"expr"`
	LegendFormat string `json:"legendFormat,omitempty"`
	Step         uint   `json:"step,omitempty"`
}

func NewPrometheusQuery(datasourceName string) *PrometheusQuery {
	return &PrometheusQuery{
		commonQuery: commonQuery{
			Datasource: &datasourceName,
		},
	}
}

func (q *PrometheusQuery) RefID() string {
	return q.commonQuery.RefID
}

func (q *PrometheusQuery) Datasource() string {
	return *q.commonQuery.Datasource
}

func (q *PrometheusQuery) commonOptions() *commonQuery {
	return &q.commonQuery
}

// GraphiteQuery is query specific options for Graphite datasource.
type GraphiteQuery struct {
	commonQuery

	Target     string `json:"target"`
	TargetFull string `json:"targetFull,omitempty"`
}

func NewGraphiteQuery(datasourceName string) *GraphiteQuery {
	return &GraphiteQuery{
		commonQuery: commonQuery{
			Datasource: &datasourceName,
		},
	}
}

func (q *GraphiteQuery) RefID() string {
	return q.commonQuery.RefID
}

func (q *GraphiteQuery) Datasource() string {
	return *q.commonQuery.Datasource
}

func (q *GraphiteQuery) commonOptions() *commonQuery {
	return &q.commonQuery
}
