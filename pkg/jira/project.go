package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// ProjectTypeClassic is a classic project type.
	ProjectTypeClassic = "classic"
	// ProjectTypeNextGen is a next gen project type.
	ProjectTypeNextGen = "next-gen"
)

// Project fetches response from /project endpoint.
func (c *Client) Project() ([]*Project, error) {
	res, err := c.GetV2(context.Background(), "/project?expand=lead", nil)
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

	var out []*Project

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}

// GetProject fetches a single project by key or id.
func (c *Client) GetProject(keyOrID string) (*Project, error) {
	res, err := c.GetV2(context.Background(), fmt.Sprintf("/project/%s?expand=lead", keyOrID), nil)
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

	var out Project
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
