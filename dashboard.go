package grafana

import (
	"encoding/json"
	"time"
)

type DashboardID uint64

type dashboardStyle string

const (
	dashboardDarkStyle  dashboardStyle = "dark"
	dashboardLightStyle dashboardStyle = "light"
)

type Dashboard struct {
	ID      DashboardID `json:"id"`
	Version uint64      `json:"version"`

	Editable      bool           `json:"editable"`
	GraphTooltip  uint8          `json:"graphTooltip"`
	HideControls  bool           `json:"hideControls"`
	Rows          []*Row         `json:"rows"`
	Style         dashboardStyle `json:"style"`
	Timezone      string         `json:"timezone"`
	Title         string         `json:"title"`
	tags          []string
	Meta          *DashboardMeta `json:"meta,omitempty"`
	SchemaVersion int            `json:"schemaVersion"`
}

// NewDashboard creates new Dashboard.
func NewDashboard(title string) *Dashboard {
	return &Dashboard{
		Title:         title,
		Editable:      true,
		SchemaVersion: 14,
		Style:         dashboardDarkStyle,
	}
}

func (d *Dashboard) String() string {
	return Stringify(d)
}

// Tags is a getter for Dashboard tags field
func (d *Dashboard) Tags() []string {
	return d.tags
}

// SetTags sets new tags to dashboard
func (d *Dashboard) SetTags(tags ...string) {
	newTags := []string{}
	uniqTags := make(map[string]bool)
	for _, tag := range tags {
		if _, ok := uniqTags[tag]; ok {
			continue
		}

		uniqTags[tag] = true
		newTags = append(newTags, tag)
	}

	d.tags = newTags
}

// AddTags adds given tags to dashboard. This method keeps uniqueness of tags.
func (d *Dashboard) AddTags(tags ...string) {
	tagFound := make(map[string]bool, len(d.tags))
	for _, tag := range d.tags {
		tagFound[tag] = true
	}

	for _, tag := range tags {
		if _, ok := tagFound[tag]; ok {
			continue
		}
		d.tags = append(d.tags, tag)
	}
}

// RemoveTags removes given tags from dashboard. Does nothing if tag is not found.
func (d *Dashboard) RemoveTags(tags ...string) {
	tagIndex := make(map[string]int, len(d.tags))
	for i, tag := range d.tags {
		tagIndex[tag] = i
	}

	for _, tag := range tags {
		if i, ok := tagIndex[tag]; ok {
			d.tags = append(d.tags[:i], d.tags[i+1:]...)
		}
	}
}

// UnmarshalJSON implements json.Unmarshaler interface
func (d *Dashboard) UnmarshalJSON(data []byte) error {
	type JSONDashboard Dashboard
	inDashboard := struct {
		*JSONDashboard
		Tags []string `json:"tags"`
	}{
		JSONDashboard: (*JSONDashboard)(d),
	}
	if err := json.Unmarshal(data, &inDashboard); err != nil {
		return err
	}

	d.tags = inDashboard.Tags

	return nil
}

// MarshalJSON implements json.Marshaler interface
func (d *Dashboard) MarshalJSON() ([]byte, error) {
	type JSONDashboard Dashboard
	dd := (*JSONDashboard)(d)
	dd.Meta = nil

	return json.Marshal(&struct {
		*JSONDashboard
		Tags []string       `json:"tags"`
		Meta *DashboardMeta `json:"-"`
	}{
		JSONDashboard: dd,
		Tags:          d.Tags(),
	})
}

type DashboardMeta struct {
	Slug    string `json:"slug"`
	Type    string `json:"type"`
	Version int    `json:"version"`

	CanEdit bool `json:"canEdit"`
	CanSave bool `json:"canSave"`
	CanStar bool `json:"canStar"`

	Created   time.Time `json:"created"`
	CreatedBy string    `json:"createdBy"`
	Expires   time.Time `json:"expires"`
	Updated   time.Time `json:"updated"`
	UpdatedBy string    `json:"updatedBy"`
}

func (dm *DashboardMeta) String() string {
	return Stringify(dm)
}

type Row struct {
	Collapsed bool         `json:"collapse"`
	Editable  bool         `json:"editable"`
	Height    string       `json:"height"`
	Panels    []*TextPanel `json:"panels"`
	RepeatFor string       `json:"repeat"` // repeat row for given variable
	ShowTitle bool         `json:"showTitle"`
	Title     string       `json:"title"`
	TitleSize string       `json:"titleSize"` // TODO: validation: h1-h6
}

// MarshalJSON implements encoding/json.Marshaler
func (r *Row) MarshalJSON() ([]byte, error) {
	for i, p := range r.Panels {
		p.id = uint(i + 1)
	}
	type JSONRow Row
	jr := (*JSONRow)(r)
	return json.Marshal(jr)
}

// NewRow creates new Row with somw defaults.
func NewRow() *Row {
	return &Row{
		Editable: true,
	}
}

type TextPanelMode string

const (
	TextPanelHTMLMode     TextPanelMode = "html"
	TextPanelMarkdownMode TextPanelMode = "markdown"
	TextPanelTextMode     TextPanelMode = "text"
)

type TextPanel struct {
	Content string        `json:"content"`
	Mode    TextPanelMode `json:"mode"`

	// General options
	id          uint
	Description string       `json:"description"`
	Height      string       `json:"height"`
	Links       []*PanelLink `json:"links"`
	MinSpan     int          `json:"minSpan"`        // TODO: valid values: 1-12
	Span        int          `json:"span,omitempty"` // TODO: valid values: 1-12
	Title       string       `json:"title"`
	Transparent bool         `json:"transparent"`
	Type        string       `json:"type"` // required
}

func (p *TextPanel) MarshalJSON() ([]byte, error) {
	type JSONPanel TextPanel
	jp := (*JSONPanel)(p)
	return json.Marshal(&struct {
		*JSONPanel
		ID uint `json:"id"`
	}{
		JSONPanel: jp,
		ID:        jp.id,
	})
}

// NewTextPanel creates new "Text" panel.
func NewTextPanel(mode TextPanelMode) *TextPanel {
	return &TextPanel{
		Mode:    mode,
		Type:    "text",
		MinSpan: 12,
	}
}

type PanelLink struct {
	IncludeVars  bool   `json:"includeVars"`
	KeepTime     bool   `json:"keepTime"`
	Params       string `json:"params"`
	OpenInNewTab bool   `json:"targetBlank"`
	Type         string `json:"type"` // TODO validation: absolute/dashboard

	// type=absolute
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`

	// type=dashboard
	DashboardURI string `json:"dashUri,omitempty"`   // TODO: validation. should be valid dashboard
	Dashboard    string `json:"dashboard,omitempty"` // actually it's title
}

// NewPanelLink creates new PanelLink
func NewPanelLink(panelType string) *PanelLink {
	return &PanelLink{
		Type: panelType, // TODO: validation
	}
}
