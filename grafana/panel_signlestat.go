package grafana

import "encoding/json"

// SinglestatPanel  represents Singlestat panel.
type SinglestatPanel struct {
	PanelGeneralOptions

	// Options. Value.
	ValueName       string `json:"valueName"`     // Stat: min/max/avg/current/total/name/first/delta/diff/range
	ValueFontSize   string `json:"valueFontSize"` // 0%-100% TODO: validation
	Postfix         string `json:"postfix"`
	PostfixFontSize string `json:"postfixFontSize"` // 0%-100% TODO: validation
	Prefix          string `json:"prefix"`
	PrefixFontSize  string `json:"prefixFontSize"` // 0%-100% TODO: validation
	Format          string `json:"format"`         // Unit option. TODO: make a custom type with constants

	// Options. Coloring.
	// Colorize background or not
	ColorBackground bool `json:"colorBackground"`
	// Colorize value or not
	ColorValue bool     `json:"colorValue"`
	Thresholds string   `json:"thresholds"` // comma separated values "x,x". TODO: validation
	Colors     []string `json:"colors"`     // array of 3 colors, ie. rgba(50, 172, 45, 0.97)

	// Options. Spark lines.
	SparkLine struct {
		Show       bool   `json:"show"`
		FullHeight bool   `json:"full"`
		LineColor  string `json:"lineColor"` // ie. rgba(50, 172, 45, 0.97)
		FillColor  string `json:"fillColor"` // ie. rgba(50, 172, 45, 0.97)
	} `json:"sparkline"`

	// Options. Gauge
	Gauge struct {
		Show             bool `json:"show"`
		MaxValue         int  `json:"maxValue"`
		MinValue         int  `json:"minValue"`
		ThresholdLabels  bool `json:"thresholdLabels"`
		ThresholdMarkers bool `json:"thresholdMarkers"`
	} `json:"gauge"`

	generalOptions PanelGeneralOptions
}

// NewSinglestatPanel creates new "Singlestat" panel.
func NewSinglestatPanel() *SinglestatPanel {
	return &SinglestatPanel{
		generalOptions: PanelGeneralOptions{
			panelType: singlestatPanel,
			MinSpan:   12,
			Links:     []*PanelLink{}, // to make [] in marshaling instead of nil
		},
	}
}

func (p *SinglestatPanel) GeneralOptions() *PanelGeneralOptions {
	return &p.generalOptions
}

func (p *SinglestatPanel) MarshalJSON() ([]byte, error) {
	type JSONPanel SinglestatPanel
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

func (p *SinglestatPanel) UnmarshalJSON(data []byte) error {
	type JSONPanel SinglestatPanel
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
