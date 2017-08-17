package grafana

import (
	"context"
)

type DashboardsService struct {
	client *Client
}

func NewDashboardsService(client *Client) *DashboardsService {
	return &DashboardsService{
		client: client,
	}
}

func (ds *DashboardsService) Search(ctx context.Context) ([]*DashboardHit, error) {
	req, err := ds.client.NewRequest(ctx, "GET", "/api/search", nil)
	if err != nil {
		return nil, err
	}

	var hits []*DashboardHit
	_, err = ds.client.Do(req, &hits)
	if err != nil {
		return nil, err
	}

	return hits, nil
}

type DashboardHit struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

func (dh *DashboardHit) String() string {
	return Stringify(dh)
}
