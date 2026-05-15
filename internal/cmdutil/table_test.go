package cmdutil

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTable(t *testing.T) {
	headers := []string{"KEY", "NAME", "TYPE"}
	rows := [][]string{
		{"FRST", "First", "classic"},
		{"SCND", "Second", "next-gen"},
	}

	cases := []struct {
		name     string
		opts     TableRenderOptions
		expected string
	}{
		{
			name:     "plain with default tab delimiter writes tab-separated rows through tabwriter",
			opts:     TableRenderOptions{Plain: true},
			expected: "KEY\tNAME\tTYPE\nFRST\tFirst\tclassic\nSCND\tSecond\tnext-gen\n",
		},
		{
			name: "plain with custom delimiter writes raw delim-separated",
			opts: TableRenderOptions{Plain: true, Delimiter: "|"},
			expected: "KEY|NAME|TYPE\n" +
				"FRST|First|classic\n" +
				"SCND|Second|next-gen\n",
		},
		{
			name: "plain + no-headers drops the header row",
			opts: TableRenderOptions{Plain: true, NoHeaders: true, Delimiter: "|"},
			expected: "FRST|First|classic\n" +
				"SCND|Second|next-gen\n",
		},
		{
			name: "csv emits RFC 4180 CSV",
			opts: TableRenderOptions{CSV: true},
			expected: "KEY,NAME,TYPE\n" +
				"FRST,First,classic\n" +
				"SCND,Second,next-gen\n",
		},
		{
			name: "csv + no-headers drops header row",
			opts: TableRenderOptions{CSV: true, NoHeaders: true},
			expected: "FRST,First,classic\n" +
				"SCND,Second,next-gen\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := captureStdout(t, func() {
				assert.NoError(t, RenderTable(headers, rows, tc.opts))
			})
			assert.Equal(t, tc.expected, out)
		})
	}
}

func TestRenderTable_CSVQuotesSpecialChars(t *testing.T) {
	headers := []string{"NAME", "BODY"}
	rows := [][]string{
		{"row1", "has,comma"},
		{"row2", "has\"quote"},
		{"row3", "has\nnewline"},
	}

	out := captureStdout(t, func() {
		assert.NoError(t, RenderTable(headers, rows, TableRenderOptions{CSV: true}))
	})

	assert.Equal(t,
		`NAME,BODY
row1,"has,comma"
row2,"has""quote"
row3,"has
newline"
`, out)
}

func TestRenderJSON(t *testing.T) {
	v := map[string]any{
		"key":   "FOO",
		"count": 2,
	}

	out := captureStdout(t, func() {
		assert.NoError(t, RenderJSON(v))
	})

	// JSON map key order is sorted by encoding/json.
	assert.Equal(t, `{
  "count": 2,
  "key": "FOO"
}
`, out)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	fn()

	assert.NoError(t, w.Close())
	os.Stdout = orig
	<-done

	return strings.ReplaceAll(buf.String(), "\r\n", "\n")
}
