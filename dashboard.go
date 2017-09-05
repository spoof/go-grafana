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

// Row is panel's row
type Row struct {
	Collapsed bool    `json:"collapse"`
	Editable  bool    `json:"editable"`
	Height    string  `json:"height"`
	Panels    []Panel `json:"panels"`
	RepeatFor string  `json:"repeat"` // repeat row for given variable
	ShowTitle bool    `json:"showTitle"`
	Title     string  `json:"title"`
	TitleSize string  `json:"titleSize"` // TODO: validation: h1-h6
}

// NewRow creates new Row with somw defaults.
func NewRow() *Row {
	return &Row{Editable: true}
}

// MarshalJSON implements encoding/json.Marshaler
func (r *Row) MarshalJSON() ([]byte, error) {
	type JSONRow Row
	jr := (*JSONRow)(r)

	for i, p := range jr.Panels {
		p.GeneralOptions().id = uint(i + 1)
	}

	return json.Marshal(jr)
}

func (r *Row) UnmarshalJSON(data []byte) error {
	type JSONRow Row
	jr := struct {
		*JSONRow
		Panels []dummyPanel `json:"panels"`
	}{
		JSONRow: (*JSONRow)(r),
	}

	if err := json.Unmarshal(data, &jr); err != nil {
		return err
	}

	panels := make([]Panel, len(jr.Panels))
	for i, p := range jr.Panels {
		panels[i] = p.Panel
	}
	r.Panels = panels
	return nil
}

type Panel interface {
	GeneralOptions() *PanelGeneralOptions
}

type panelType string

const (
	textPanel       panelType = "text"
	singlestatPanel panelType = "singlestat"
)

type dummyPanel struct {
	Type panelType `json:"type"`
	Panel
}

func (p *dummyPanel) UnmarshalJSON(data []byte) error {
	type JSONPanel dummyPanel
	var jp JSONPanel
	if err := json.Unmarshal(data, &jp); err != nil {
		return err
	}

	var panel Panel
	switch jp.Type {
	case textPanel:
		panel = new(TextPanel)
	case singlestatPanel:
		panel = new(SinglestatPanel)
	default:
		return nil
	}

	if err := json.Unmarshal(data, panel); err != nil {
		return err
	}

	p.Panel = panel
	return nil
}

type PanelGeneralOptions struct {
	id        uint
	panelType panelType

	Description string       `json:"description"`
	Height      string       `json:"height"`
	Links       []*PanelLink `json:"links"`
	MinSpan     uint         `json:"minSpan"` // TODO: valid values: 1-12
	Span        uint         `json:"span"`    // TODO: valid values: 1-12
	Title       string       `json:"title"`
	Transparent bool         `json:"transparent"`
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
