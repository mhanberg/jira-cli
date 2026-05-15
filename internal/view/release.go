package view

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

// ProjectVersionOptions is a functional option to wrap project version properties.
type ProjectVersionOptions func(*Release)

// Release is a release view.
type Release struct {
	data      []*jira.ProjectVersion
	writer    io.Writer
	buf       *bytes.Buffer
	plain     bool
	noHeaders bool
	csv       bool
	delimiter string
}

// NewRelease constructs a project release command.
func NewRelease(data []*jira.ProjectVersion, opts ...ProjectVersionOptions) *Release {
	r := Release{
		data: data,
		buf:  new(bytes.Buffer),
	}
	r.writer = tabwriter.NewWriter(r.buf, 0, tabWidth, 1, '\t', 0)

	for _, opt := range opts {
		opt(&r)
	}
	return &r
}

// WithReleaseWriter sets a writer for the project release.
func WithReleaseWriter(w io.Writer) ProjectVersionOptions {
	return func(r *Release) {
		r.writer = w
	}
}

// WithReleasePlain toggles plain output (no interactive pager).
func WithReleasePlain(v bool) ProjectVersionOptions {
	return func(r *Release) {
		r.plain = v
	}
}

// WithReleaseNoHeaders skips the header row.
func WithReleaseNoHeaders(v bool) ProjectVersionOptions {
	return func(r *Release) {
		r.noHeaders = v
	}
}

// WithReleaseCSV emits CSV output.
func WithReleaseCSV(v bool) ProjectVersionOptions {
	return func(r *Release) {
		r.csv = v
	}
}

// WithReleaseDelimiter sets the column delimiter in plain mode. Default is "\t".
func WithReleaseDelimiter(d string) ProjectVersionOptions {
	return func(r *Release) {
		r.delimiter = d
	}
}

// Render renders the project release view.
func (r Release) Render() error {
	data := r.tableData()

	if r.csv {
		return renderCSV(r.outputWriter(), data)
	}

	delim := r.delimiter
	if delim == "" {
		delim = "\t"
	}
	if r.plain && delim != "\t" {
		return renderPlain(r.outputWriter(), data, delim)
	}

	for _, row := range data {
		for i, cell := range row {
			_, _ = fmt.Fprint(r.writer, cell)
			if i != len(row)-1 {
				_, _ = fmt.Fprint(r.writer, "\t")
			}
		}
		_, _ = fmt.Fprintln(r.writer)
	}
	if tw, ok := r.writer.(*tabwriter.Writer); ok {
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	if _, ok := r.writer.(*tabwriter.Writer); !ok {
		return nil
	}
	if r.plain {
		_, err := fmt.Fprint(os.Stdout, r.buf.String())
		return err
	}
	return tui.PagerOut(r.buf.String())
}

func (r Release) tableData() tui.TableData {
	var data tui.TableData
	if !r.noHeaders {
		data = append(data, r.header())
	}
	for _, d := range r.data {
		desc := ""
		if d.Description != nil {
			desc = fmt.Sprint(d.Description)
		}
		data = append(data, []string{
			d.ID,
			prepareTitle(d.Name),
			fmt.Sprintf("%v", d.Released),
			desc,
		})
	}
	return data
}

func (r Release) outputWriter() io.Writer {
	if _, isTab := r.writer.(*tabwriter.Writer); isTab {
		return os.Stdout
	}
	return r.writer
}

func (r Release) header() []string {
	return []string{
		"ID",
		"NAME",
		"RELEASED",
		"DESCRIPTION",
	}
}

func (r Release) printHeader() {
	headers := r.header()
	end := len(headers) - 1
	for i, h := range headers {
		_, _ = fmt.Fprintf(r.writer, "%s", h)
		if i != end {
			_, _ = fmt.Fprintf(r.writer, "\t")
		}
	}
	_, _ = fmt.Fprintln(r.writer)
}
