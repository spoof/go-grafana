package grafana

import (
	"context"
)

type DashboardID uint64

type Dashboard struct {
	ID           DashboardID `json:"id"`
	Editable     bool        `json:"editable"`
	GraphTooltip uint8       `json:"graphTooltip"`
	HideControls bool        `json:"hideControls"`
	Rows         []*Row      `json:"rows"`
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

type DashboardsService struct {
	client *Client
}

func NewDashboardsService(client *Client) *DashboardsService {
	return &DashboardsService{
		client: client,
	}
}

func (ds *DashboardsService) Search(ctx context.Context) ([]*Dashboard, error) {
	req, err := ds.client.NewRequest(ctx, "GET", "/api/search", nil)
	if err != nil {
		return nil, err
	}

	var dashboards []*Dashboard
	_, err = ds.client.Do(req, &dashboards)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
