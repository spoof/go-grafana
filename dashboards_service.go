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

type DashboardSearchOptions struct {
	Query     string   `url:"query,omitempty"`
	Tags      []string `url:"tags,omitempty"`
	IsStarred bool     `url:"starred,omitempty"`
	Limit     int      `url:"limit,omitempty"`
}

func (ds *DashboardsService) Search(ctx context.Context, opt *DashboardSearchOptions) ([]*DashboardHit, error) {
	u := "/api/search"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, err
	}

	req, err := ds.client.NewRequest(ctx, "GET", u, nil)
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
