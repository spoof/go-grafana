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

import "encoding/json"

type (
	// DatasourceID represents id type of datasource
	DatasourceID uint
	// httpAccessType represents type of http access options
	httpAccessType string
)

// Types of HTTP access to datasource
const (
	HTTPAccesProxy  httpAccessType = "proxy"
	HTTPAccesDirect httpAccessType = "direct"
)

// datasourceType is type of datasource
type datasourceType string

// Types of datasource
const (
	GraphiteDatasource   datasourceType = "graphite"
	PrometheusDatasource datasourceType = "prometheus"
)

// Datasource represents datasource entity of Grafana.
type Datasource struct {
	id    DatasourceID
	OrgID OrgID `json:"orgId"`

	Name              string         `json:"name"`
	Type              datasourceType `json:"type"`
	Access            httpAccessType `json:"access"`
	URL               string         `json:"url"`
	Password          string         `json:"password"`
	User              string         `json:"user"`
	Database          string         `json:"database"`
	BasicAuth         bool           `json:"basicAuth"`
	BasicAuthUser     string         `json:"basicAuthUser"`
	BasicAuthPassword string         `json:"basicAuthPassword"`
	IsDefault         bool           `json:"isDefault"`
	WithCredentials   bool           `json:"withCredentials"`

	// SecureJsonData    securejsondata.SecureJsonData
	// JsonData          *simplejson.Json
	// "typeLogoUrl": "",

}

// ID returns id of Datasource
func (d *Datasource) ID() DatasourceID {
	return d.id
}

// MarshalJSON implements json.Marshaler interface
func (d *Datasource) MarshalJSON() ([]byte, error) {
	type JSONDatasource Datasource
	jd := struct {
		*JSONDatasource
		ID DatasourceID `json:"id"`
	}{
		JSONDatasource: (*JSONDatasource)(d),
		ID:             d.id,
	}
	return json.Marshal(jd)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Datasource) UnmarshalJSON(data []byte) error {
	type JSONDatasource Datasource
	jd := struct {
		*JSONDatasource
		ID *DatasourceID `json:"id"`
	}{
		JSONDatasource: (*JSONDatasource)(d),
		ID:             &d.id,
	}

	return json.Unmarshal(data, &jd)
}
