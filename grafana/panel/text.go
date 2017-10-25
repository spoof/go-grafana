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

	generalOptions GeneralOptions
}

// NewTextPanel creates new "Text" panel.
func NewTextPanel(mode TextPanelMode) *TextPanel {
	return &TextPanel{
		Mode: mode,
	}
}

func (p *TextPanel) GeneralOptions() *GeneralOptions {
	return &p.generalOptions
}
