package view

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestProjectRender(t *testing.T) {
	var b bytes.Buffer

	//nolint:unused
	type lead struct {
		Name string `json:"displayName"`
	}

	data := []*jira.Project{
		{Key: "FRST", Name: "First", Lead: lead{Name: "Person A"}, Type: jira.ProjectTypeClassic},
		{Key: "SCND", Name: "[2] Second", Lead: lead{Name: "Person B"}, Type: jira.ProjectTypeNextGen},
		{Key: "THIRD", Name: "Third", Lead: lead{Name: "Person C"}, Type: jira.ProjectTypeClassic},
	}
	project := NewProject(data, WithProjectWriter(&b))
	assert.NoError(t, project.Render())

	expected := `KEY	NAME	TYPE	LEAD
FRST	First	classic	Person A
SCND	[2[] Second	next-gen	Person B
THIRD	Third	classic	Person C
`
	assert.Equal(t, expected, b.String())
}

func TestProjectRenderModes(t *testing.T) {
	//nolint:unused
	type lead struct {
		Name string `json:"displayName"`
	}

	data := []*jira.Project{
		{Key: "FRST", Name: "First", Lead: lead{Name: "Person A"}, Type: jira.ProjectTypeClassic},
		{Key: "SCND", Name: "Second", Lead: lead{Name: "Person B"}, Type: jira.ProjectTypeNextGen},
	}

	t.Run("no headers", func(t *testing.T) {
		var b bytes.Buffer
		v := NewProject(data, WithProjectWriter(&b), WithProjectNoHeaders(true))
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"FRST\tFirst\tclassic\tPerson A\nSCND\tSecond\tnext-gen\tPerson B\n",
			b.String())
	})

	t.Run("csv", func(t *testing.T) {
		var b bytes.Buffer
		v := NewProject(data, WithProjectWriter(&b), WithProjectCSV(true))
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"KEY,NAME,TYPE,LEAD\nFRST,First,classic,Person A\nSCND,Second,next-gen,Person B\n",
			b.String())
	})

	t.Run("plain with custom delimiter", func(t *testing.T) {
		var b bytes.Buffer
		v := NewProject(data,
			WithProjectWriter(&b),
			WithProjectPlain(true),
			WithProjectDelimiter("|"),
		)
		assert.NoError(t, v.Render())
		assert.Equal(t,
			"KEY|NAME|TYPE|LEAD\nFRST|First|classic|Person A\nSCND|Second|next-gen|Person B\n",
			b.String())
	})
}
