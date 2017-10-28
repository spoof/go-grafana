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
	"github.com/spoof/go-grafana/pkg/field"
)

type GeneralOptions struct {
	Description string            `json:"description"`
	Height      field.ForceString `json:"height"`
	Links       []PanelLink       `json:"links"`
	MinSpan     uint              `json:"minSpan"` // TODO: valid values: 1-12
	Span        uint              `json:"span"`    // TODO: valid values: 1-12
	Title       string            `json:"title"`
	Transparent bool              `json:"transparent"`
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

// Query is interface describes behaviour that we want from panel's queries
type Query interface {
	Datasource() string
}
