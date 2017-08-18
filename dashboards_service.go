package grafana

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// DashboardsService communicates with dashboard methods of the Grafana API.
type DashboardsService struct {
	client *Client
}

// NewDashboardsService returns a new DashboardsService.
func NewDashboardsService(client *Client) *DashboardsService {
	return &DashboardsService{
		client: client,
	}
}

// ErrDashboardNotFound represents an error if dashboard not found.
var ErrDashboardNotFound = errors.New("Dashboard not found")

// Get fetches a dashboard by given slug.
//
// Grafana API docs: http://docs.grafana.org/http_api/dashboard/#get-dashboard
func (ds *DashboardsService) Get(ctx context.Context, slug string) (*Dashboard, error) {
	u := fmt.Sprintf("/api/dashboards/db/%s", slug)
	req, err := ds.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var d Dashboard
	if resp, err := ds.client.Do(req, &d); err != nil {
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrDashboardNotFound
			}
		}
		return nil, err
	}

	return &d, nil
}

// DashboardSearchOptions specifies the optional parameters to the
// DashboardsService.Search method.
type DashboardSearchOptions struct {
	Query     string   `url:"query,omitempty"`
	Tags      []string `url:"tags,omitempty"`
	IsStarred bool     `url:"starred,omitempty"`
	Limit     int      `url:"limit,omitempty"`
}

// Search searches dashboards with given criteria
//
//  Grafana API docs: http://docs.grafana.org/http_api/dashboard/#search-dashboards
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

// DashboardHit represents a found by DashboardsService.Search dashboard
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
