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

// GraphPanel represents Graph panel
type GraphPanel struct {
	generalOptions GeneralOptions
	//queriesOptions QueriesOptions
}

// NewGraphPanel creates new Graph panel.
func NewGraphPanel() *GraphPanel {
	return &GraphPanel{}
}

func (p *GraphPanel) GeneralOptions() *GeneralOptions {
	return &p.generalOptions
}
