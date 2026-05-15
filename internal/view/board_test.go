package view

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

func TestBoardRender(t *testing.T) {
	var b bytes.Buffer

	data := []*jira.Board{
		{ID: 1, Name: "First", Type: "scrum"},
		{ID: 2, Name: "[2] Second", Type: "kanban"},
		{ID: 3, Name: "Third", Type: "nextgen"},
	}
	board := NewBoard(data, WithBoardWriter(&b))
	assert.NoError(t, board.Render())

	expected := `ID	NAME	TYPE
1	First	scrum
2	[2[] Second	kanban
3	Third	nextgen
`
	assert.Equal(t, expected, b.String())
}

func TestBoardRenderModes(t *testing.T) {
	data := []*jira.Board{
		{ID: 1, Name: "First", Type: "scrum"},
		{ID: 2, Name: "Second", Type: "kanban"},
	}

	t.Run("no headers", func(t *testing.T) {
		var b bytes.Buffer
		v := NewBoard(data, WithBoardWriter(&b), WithBoardNoHeaders(true))
		assert.NoError(t, v.Render())
		assert.Equal(t, "1\tFirst\tscrum\n2\tSecond\tkanban\n", b.String())
	})

	t.Run("csv", func(t *testing.T) {
		var b bytes.Buffer
		v := NewBoard(data, WithBoardWriter(&b), WithBoardCSV(true))
		assert.NoError(t, v.Render())
		assert.Equal(t, "ID,NAME,TYPE\n1,First,scrum\n2,Second,kanban\n", b.String())
	})

	t.Run("plain with custom delimiter", func(t *testing.T) {
		var b bytes.Buffer
		v := NewBoard(data,
			WithBoardWriter(&b),
			WithBoardPlain(true),
			WithBoardDelimiter("|"),
		)
		assert.NoError(t, v.Render())
		assert.Equal(t, "ID|NAME|TYPE\n1|First|scrum\n2|Second|kanban\n", b.String())
	})
}
