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
	tags         []string
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
		uniqTags[tag] = true
	}

	for tag := range uniqTags {
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

type Row struct {
	Collapse bool     `json:"collapse"`
	Editable bool     `json:"editable"`
	Height   string   `json:"height"`
	Title    string   `json:"title"`
	Panels   []*Panel `json:"panels"`
}

type Panel struct {
}
