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

// DatasourcesService communicates with datasource methods of the Grafana API.
type DatasourcesService struct {
	client *Client
}

// NewDatasourcesService returns a new DatasourcesService.
func NewDatasourcesService(client *Client) *DatasourcesService {
	return &DatasourcesService{
		client: client,
	}
}

// GetAll fetches all datasources.
//
// Grafana API docs: http://docs.grafana.org/http_api/data_source/#get-all-datasources
func (s *DatasourcesService) GetAll(ctx context.Context) ([]*grafana.Datasource, error) {
	u := "/api/datasources"
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var datasources []*grafana.Datasource
	if _, err := s.client.Do(req, &datasources); err != nil {
		return nil, err
	}

	return datasources, nil
}

// ErrDatasourceNotFound represents an error if datasource not found.
var ErrDatasourceNotFound = errors.New("Datasource not found")

// GetByID fetches datasource by given id.
//
// Grafana API docs: http://docs.grafana.org/http_api/data_source/#get-a-single-data-sources-by-id
func (s *DatasourcesService) GetByID(ctx context.Context, id grafana.DatasourceID) (*grafana.Datasource, error) {
	u := fmt.Sprintf("/api/datasources/%d", id)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var d grafana.Datasource
	if resp, err := s.client.Do(req, &d); err != nil {
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrDatasourceNotFound
			}
		}

		return nil, err
	}

	return &d, nil
}

// GetByName fetches datasource with given name.
//
// Grafana API docs: http://docs.grafana.org/http_api/data_source/#get-a-single-data-source-by-name
func (s *DatasourcesService) GetByName(ctx context.Context, name string) (*grafana.Datasource, error) {
	if name == "" {
		return nil, errors.New("Name cannot be empty")
	}

	u := fmt.Sprintf("/api/datasources/name/%s", name)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, err
	}

	var d grafana.Datasource
	if resp, err := s.client.Do(req, &d); err != nil {
		if resp != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrDatasourceNotFound
			}
		}

		return nil, err
	}

	return &d, nil
}
