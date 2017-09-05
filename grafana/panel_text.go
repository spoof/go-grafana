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

// TextPanelMode is a type of Text panel.
type TextPanelMode string

// This is all possible types (modes) of Text panel.
const (
	TextPanelHTMLMode     TextPanelMode = "html"
	TextPanelMarkdownMode TextPanelMode = "markdown"
	TextPanelTextMode     TextPanelMode = "text"
)

// TextPanel represents Text Panel
type TextPanel struct {
	Content string        `json:"content"`
	Mode    TextPanelMode `json:"mode"`

	generalOptions PanelGeneralOptions
}

// NewTextPanel creates new "Text" panel.
func NewTextPanel(mode TextPanelMode) *TextPanel {
	return &TextPanel{
		Mode: mode,
		generalOptions: PanelGeneralOptions{
			panelType: textPanel,
			MinSpan:   12,
		},
	}
}

func (p *TextPanel) GeneralOptions() *PanelGeneralOptions {
	return &p.generalOptions
}

func (p *TextPanel) MarshalJSON() ([]byte, error) {
	type JSONPanel TextPanel
	jp := struct {
		*JSONPanel
		*PanelGeneralOptions
		ID   uint      `json:"id"`
		Type panelType `json:"type"`
	}{
		JSONPanel:           (*JSONPanel)(p),
		PanelGeneralOptions: p.GeneralOptions(),
		ID:                  p.GeneralOptions().id,
		Type:                p.GeneralOptions().panelType,
	}
	return json.Marshal(jp)
}

func (p *TextPanel) UnmarshalJSON(data []byte) error {
	type JSONPanel TextPanel
	jp := struct {
		*JSONPanel
		*PanelGeneralOptions
		Type *panelType `json:"type"`
	}{
		JSONPanel:           (*JSONPanel)(p),
		PanelGeneralOptions: p.GeneralOptions(),
		Type:                &p.GeneralOptions().panelType,
	}

	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}

	return nil
}
