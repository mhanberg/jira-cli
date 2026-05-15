package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestProxySearchAll_V3_TokenPagination(t *testing.T) {
	pages := []string{
		`{"isLast": false, "nextPageToken": "tok-2", "issues": [{"key":"P-1"},{"key":"P-2"}]}`,
		`{"isLast": false, "nextPageToken": "tok-3", "issues": [{"key":"P-3"},{"key":"P-4"}]}`,
		`{"isLast": true,  "nextPageToken": "",      "issues": [{"key":"P-5"}]}`,
	}
	tokenToPage := map[string]int{"": 0, "tok-2": 1, "tok-3": 2}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/3/search/jql", r.URL.Path)
		token := r.URL.Query().Get("nextPageToken")
		idx, ok := tokenToPage[token]
		assert.True(t, ok, "unexpected nextPageToken %q", token)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(pages[idx]))
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("installation", jira.InstallationTypeCloud)
	t.Cleanup(viper.Reset)

	client := jira.NewClient(jira.Config{Server: server.URL}, jira.WithTimeout(3*time.Second))

	out, err := ProxySearchAll(client, "project=TEST", 50)
	assert.NoError(t, err)

	keys := keysOf(out.Issues)
	assert.Equal(t, []string{"P-1", "P-2", "P-3", "P-4", "P-5"}, keys)
	assert.True(t, out.IsLast)
}

func TestProxySearchAll_V2_OffsetPagination(t *testing.T) {
	// Page 0: 3 issues. Page 1: 3 issues. Page 2: 1 issue (shorter than page size → last).
	pageRanges := map[string]string{
		"0": `{"issues": [{"key":"P-1"},{"key":"P-2"},{"key":"P-3"}]}`,
		"3": `{"issues": [{"key":"P-4"},{"key":"P-5"},{"key":"P-6"}]}`,
		"6": `{"issues": [{"key":"P-7"}]}`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/search", r.URL.Path)
		startAt := r.URL.Query().Get("startAt")
		body, ok := pageRanges[startAt]
		assert.True(t, ok, "unexpected startAt %q", startAt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	viper.Reset()
	viper.Set("installation", jira.InstallationTypeLocal)
	t.Cleanup(viper.Reset)

	client := jira.NewClient(jira.Config{Server: server.URL}, jira.WithTimeout(3*time.Second))

	out, err := ProxySearchAll(client, "project=TEST", 3)
	assert.NoError(t, err)

	keys := keysOf(out.Issues)
	assert.Equal(t, []string{"P-1", "P-2", "P-3", "P-4", "P-5", "P-6", "P-7"}, keys)
	assert.True(t, out.IsLast)
}

func TestProxyBoardsAll_PaginatesUntilShortPage(t *testing.T) {
	pageBodies := map[string]string{
		"0":   boardsPageBody(50, 0),
		"50":  boardsPageBody(50, 50),
		"100": boardsPageBody(7, 100),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/agile/1.0/board", r.URL.Path)
		startAt := r.URL.Query().Get("startAt")
		body, ok := pageBodies[startAt]
		assert.True(t, ok, "unexpected startAt %q", startAt)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	viper.Reset()
	t.Cleanup(viper.Reset)

	client := jira.NewClient(jira.Config{Server: server.URL}, jira.WithTimeout(3*time.Second))

	out, err := ProxyBoardsAll(client, "TEST", jira.BoardTypeAll)
	assert.NoError(t, err)
	assert.Len(t, out.Boards, 107)
}

// boardsPageBody emits a JSON document mimicking the agile /board response
// shape: n boards starting at id=startAt+1.
func boardsPageBody(n, startAt int) string {
	body := `{"values":[`
	for i := 0; i < n; i++ {
		if i > 0 {
			body += ","
		}
		body += fmt.Sprintf(`{"id":%d,"name":"Board %d","type":"scrum"}`, startAt+i+1, startAt+i+1)
	}
	body += `]}`
	return body
}

func keysOf(issues []*jira.Issue) []string {
	out := make([]string, len(issues))
	for i, iss := range issues {
		out[i] = iss.Key
	}
	return out
}
