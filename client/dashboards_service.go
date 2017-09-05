// Copyright 2017 Sergey Safonov
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spoof/go-grafana/grafana"
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
func (ds *DashboardsService) Get(ctx context.Context, slug string) (*grafana.Dashboard, error) {
	u := fmt.Sprintf("/api/dashboards/db/%s", slug)
	req, err := ds.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var dResp dashboardGetResponse
	if resp, err := ds.client.Do(req, &dResp); err != nil {
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrDashboardNotFound
			}
		}
		return nil, err
	}

	d := dResp.Dashboard
	d.Meta = dResp.Meta
	return &d, nil
}

type dashboardGetResponse struct {
	Dashboard grafana.Dashboard      `json:"dashboard"`
	Meta      *grafana.DashboardMeta `json:"meta"`
}

// Save creates a new dashboard or updates existing one.
//
// Grafana API docs: http://docs.grafana.org/http_api/dashboard/#create-update-dashboard
func (ds *DashboardsService) Save(ctx context.Context, dashboard *grafana.Dashboard, overwrite bool) error {
	u := "/api/dashboards/db"

	dReq := dashboardCreateRequest{dashboard, overwrite}
	req, err := ds.client.NewRequest(ctx, "POST", u, dReq)
	if err != nil {
		return err
	}

	var respBody struct {
		Slug    string `json:"slug"`
		Status  string `json:"status"`
		Version int    `json:"version"`
	}
	if _, err := ds.client.Do(req, &respBody); err != nil {
		// TODO: handle errors properly
		// 400 {"message":"Dashboard title cannot be empty", "error": ...}
		// 404 {"status": "not-found", "message": err.Error()}
		// 412 {"status": "name-exists", "message": err.Error()}
		// 412 {"status": "version-mismatch", "message": err.Error()}
		// 412 {"status": "plugin-dashboard", "message": message}
		// 500 {"message": "failed to get quota", "error": ...}

		return err
	}

	// To make our dashboard in sync with Grafana's one
	// we need to refetch just saved dashboard by using `Get dashboard` API.
	if respBody.Status == "success" {
		d, err := ds.Get(ctx, respBody.Slug)
		if err != nil {
			return err
		}
		*dashboard = *d
	}

	return nil
}

type dashboardCreateRequest struct {
	Dashboard *grafana.Dashboard `json:"dashboard"`
	Overwrite bool               `json:"overwrite"`
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
	return grafana.Stringify(dh)
}
