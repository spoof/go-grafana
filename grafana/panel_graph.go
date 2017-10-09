package grafana

import "encoding/json"

// GraphPanel represents Graph panel
type GraphPanel struct {
	generalOptions PanelGeneralOptions
	queriesOptions QueriesOptions
}

// NewGraphPanel creates new Graph panel.
func NewGraphPanel() *GraphPanel {
	return &GraphPanel{
		generalOptions: PanelGeneralOptions{
			panelType: graphPanel,
			MinSpan:   12,
		},
	}
}

func (p *GraphPanel) GeneralOptions() *PanelGeneralOptions {
	return &p.generalOptions
}

func (p *GraphPanel) QueriesOptions() *QueriesOptions {
	return &p.queriesOptions
}

func (p *GraphPanel) MarshalJSON() ([]byte, error) {
	type JSONPanel GraphPanel
	jp := struct {
		*JSONPanel
		*PanelGeneralOptions
		*QueriesOptions
		ID   uint      `json:"id"`
		Type panelType `json:"type"`
	}{
		JSONPanel:           (*JSONPanel)(p),
		PanelGeneralOptions: p.GeneralOptions(),
		QueriesOptions:      p.QueriesOptions(),
		ID:                  p.GeneralOptions().id,
		Type:                p.GeneralOptions().panelType,
	}
	return json.Marshal(jp)
}

func (p *GraphPanel) UnmarshalJSON(data []byte) error {
	type JSONPanel GraphPanel
	jp := struct {
		*JSONPanel
		*PanelGeneralOptions
		*QueriesOptions
		Type *panelType `json:"type"`
	}{
		JSONPanel:           (*JSONPanel)(p),
		PanelGeneralOptions: p.GeneralOptions(),
		QueriesOptions:      p.QueriesOptions(),
		Type:                &p.GeneralOptions().panelType,
	}

	return json.Unmarshal(data, &jp)
}
