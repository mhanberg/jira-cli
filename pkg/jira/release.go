package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Release fetches response from /project/{projectIdOrKey}/version endpoint.
func (c *Client) Release(project string) ([]*ProjectVersion, error) {
	path := fmt.Sprintf("/project/%s/versions", project)
	res, err := c.Get(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out []*ProjectVersion

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}

// GetVersion fetches a single project version (release) by id.
func (c *Client) GetVersion(id string) (*ProjectVersion, error) {
	res, err := c.Get(context.Background(), fmt.Sprintf("/version/%s", id), nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out ProjectVersion
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
