package view

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestReleaseRender(t *testing.T) {
	var b bytes.Buffer

	//nolint:unused
	type lead struct {
		Name string `json:"displayName"`
	}

	data := []*jira.ProjectVersion{
		{ID: "1000", Name: "First", Released: true, Description: "Release A"},
		{ID: "1001", Name: "Second", Released: false, Description: "Release B"},
		{ID: "1002", Name: "Third", Released: false, Description: "Release C"},
	}
	project := NewRelease(data, WithReleaseWriter(&b))
	assert.NoError(t, project.Render())

	expected := `ID	NAME	RELEASED	DESCRIPTION
1000	First	true	Release A
1001	Second	false	Release B
1002	Third	false	Release C
`
	assert.Equal(t, expected, b.String())
}

func TestReleaseRenderModes(t *testing.T) {
	data := []*jira.ProjectVersion{
		{ID: "1000", Name: "First", Released: true, Description: "Release A"},
		{ID: "1001", Name: "Second", Released: false, Description: "Release B"},
	}

	t.Run("no headers", func(t *testing.T) {
		var b bytes.Buffer
		v := NewRelease(data,
			WithReleaseWriter(&b),
			WithReleaseNoHeaders(true),
		)
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"1000\tFirst\ttrue\tRelease A\n1001\tSecond\tfalse\tRelease B\n",
			b.String())
	})

	t.Run("csv", func(t *testing.T) {
		var b bytes.Buffer
		v := NewRelease(data,
			WithReleaseWriter(&b),
			WithReleaseCSV(true),
		)
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"ID,NAME,RELEASED,DESCRIPTION\n1000,First,true,Release A\n1001,Second,false,Release B\n",
			b.String())
	})

	t.Run("plain with custom delimiter", func(t *testing.T) {
		var b bytes.Buffer
		v := NewRelease(data,
			WithReleaseWriter(&b),
			WithReleasePlain(true),
			WithReleaseDelimiter("|"),
		)
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"ID|NAME|RELEASED|DESCRIPTION\n1000|First|true|Release A\n1001|Second|false|Release B\n",
			b.String())
	})
}
