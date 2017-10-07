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

type (
	DatasourceID   uint
	httpAccessType string
)

const (
	HTTPAccesProxy  = "proxy"
	HTTPAccesDirect = "direct"
)

type Datasource struct {
	ID    DatasourceID `json:"id"`
	OrgID OrgID        `json:"orgId"`

	Name              string         `json:"name"`
	Type              string         `json:"type"`
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

func (d *Datasource) String() string {
	return Stringify(d)
}
