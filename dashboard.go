package grafana

type DashboardID uint64

type Dashboard struct {
	ID           DashboardID `json:"id"`
	Editable     bool        `json:"editable"`
	GraphTooltip uint8       `json:"graphTooltip"`
	HideControls bool        `json:"hideControls"`
	Rows         []*Row      `json:"rows"`
	Slug         string      `json:"slug"`
	Style        string      `json:"style"`
	Timezone     string      `json:"timezone"`
	Title        string      `json:"title"`
	Tags         *[]string   `json:"tags"`
}

type Row struct {
	Collapse bool     `json:"collapse"`
	Editable bool     `json:"editable"`
	Height   string   `json:"height"`
	Title    string   `json:"title"`
	Panels   []*Panel `json:"panels"`
}

type Panel struct {
}
