package cmdutil

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// TableRenderOptions controls how a tabular result is emitted by RenderTable.
//
// Plain: write directly to stdout instead of piping through the configured pager.
// NoHeaders: skip the header row.
// CSV: emit RFC 4180 CSV instead of tab-aligned columns.
// Delimiter: column separator in plain mode (default "\t"). Ignored unless Plain is true.
type TableRenderOptions struct {
	Plain     bool
	NoHeaders bool
	CSV       bool
	Delimiter string
}

// RenderTable emits a header + rows respecting Plain/CSV/Delimiter modes.
// The default mode pages through the configured pager via tui.PagerOut.
func RenderTable(headers []string, rows [][]string, opts TableRenderOptions) error {
	var data [][]string
	if !opts.NoHeaders {
		data = append(data, headers)
	}
	data = append(data, rows...)

	if opts.CSV {
		return csvWrite(os.Stdout, data)
	}

	delim := opts.Delimiter
	if delim == "" {
		delim = "\t"
	}
	if opts.Plain && delim != "\t" {
		return plainWrite(os.Stdout, data, delim)
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 8, 1, '\t', 0)
	for _, row := range data {
		for i, c := range row {
			_, _ = fmt.Fprint(w, c)
			if i != len(row)-1 {
				_, _ = fmt.Fprint(w, "\t")
			}
		}
		_, _ = fmt.Fprintln(w)
	}
	if err := w.Flush(); err != nil {
		return err
	}

	if opts.Plain {
		_, err := fmt.Fprint(os.Stdout, buf.String())
		return err
	}
	return tui.PagerOut(buf.String())
}

// RenderJSON marshals v with two-space indent and prints to stdout.
func RenderJSON(v interface{}) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Println(string(out))
	return err
}

func csvWrite(w io.Writer, data [][]string) error {
	wr := csv.NewWriter(w)
	for _, row := range data {
		if err := wr.Write(row); err != nil {
			return err
		}
	}
	wr.Flush()
	return wr.Error()
}

func plainWrite(w io.Writer, data [][]string, delim string) error {
	for _, row := range data {
		for i, c := range row {
			_, _ = fmt.Fprint(w, c)
			if i != len(row)-1 {
				_, _ = fmt.Fprint(w, delim)
			}
		}
		_, _ = fmt.Fprintln(w)
	}
	return nil
}
