package grafana

import (
	"encoding/json"
	"time"
)

type DashboardID uint64

type Dashboard struct {
	ID      DashboardID `json:"id"`
	Version uint64      `json:"version"`

	Editable     bool   `json:"editable"`
	GraphTooltip uint8  `json:"graphTooltip"`
	HideControls bool   `json:"hideControls"`
	Rows         []*Row `json:"rows"`
	Style        string `json:"style"`
	Timezone     string `json:"timezone"`
	Title        string `json:"title"`
	tags         []string
	Meta         *DashboardMeta `json:"meta,omitempty"`
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
		Tags []string `json:"tags,omitempty"`
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
	Collapse bool     `json:"collapse"`
	Editable bool     `json:"editable"`
	Height   string   `json:"height"`
	Title    string   `json:"title"`
	Panels   []*Panel `json:"panels"`
}

type Panel struct {
}
