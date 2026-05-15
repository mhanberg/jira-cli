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

// ProjectOption is a functional option to wrap project properties.
type ProjectOption func(*Project)

// Project is a project view.
type Project struct {
	data      []*jira.Project
	writer    io.Writer
	buf       *bytes.Buffer
	plain     bool
	noHeaders bool
	csv       bool
	delimiter string
}

// NewProject initializes a project.
func NewProject(data []*jira.Project, opts ...ProjectOption) *Project {
	p := Project{
		data: data,
		buf:  new(bytes.Buffer),
	}
	p.writer = tabwriter.NewWriter(p.buf, 0, tabWidth, 1, '\t', 0)

	for _, opt := range opts {
		opt(&p)
	}
	return &p
}

// WithProjectWriter sets a writer for the project.
func WithProjectWriter(w io.Writer) ProjectOption {
	return func(p *Project) {
		p.writer = w
	}
}

// WithProjectPlain toggles plain output (no interactive pager).
func WithProjectPlain(v bool) ProjectOption {
	return func(p *Project) {
		p.plain = v
	}
}

// WithProjectNoHeaders skips the header row.
func WithProjectNoHeaders(v bool) ProjectOption {
	return func(p *Project) {
		p.noHeaders = v
	}
}

// WithProjectCSV emits CSV output.
func WithProjectCSV(v bool) ProjectOption {
	return func(p *Project) {
		p.csv = v
	}
}

// WithProjectDelimiter sets the column delimiter in plain mode. Default is "\t".
func WithProjectDelimiter(d string) ProjectOption {
	return func(p *Project) {
		p.delimiter = d
	}
}

// Render renders the project view.
func (p Project) Render() error {
	data := p.tableData()

	if p.csv {
		return renderCSV(p.outputWriter(), data)
	}

	delim := p.delimiter
	if delim == "" {
		delim = "\t"
	}
	if p.plain && delim != "\t" {
		return renderPlain(p.outputWriter(), data, delim)
	}

	for _, row := range data {
		for i, cell := range row {
			_, _ = fmt.Fprint(p.writer, cell)
			if i != len(row)-1 {
				_, _ = fmt.Fprint(p.writer, "\t")
			}
		}
		_, _ = fmt.Fprintln(p.writer)
	}
	if tw, ok := p.writer.(*tabwriter.Writer); ok {
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	if _, ok := p.writer.(*tabwriter.Writer); !ok {
		// Custom writer (used in tests) — content already there.
		return nil
	}
	if p.plain {
		_, err := fmt.Fprint(os.Stdout, p.buf.String())
		return err
	}
	return tui.PagerOut(p.buf.String())
}

func (p Project) tableData() tui.TableData {
	var data tui.TableData
	if !p.noHeaders {
		data = append(data, p.header())
	}
	for _, d := range p.data {
		data = append(data, []string{d.Key, prepareTitle(d.Name), d.Type, d.Lead.Name})
	}
	return data
}

// outputWriter returns where csv/non-tab-delimited output should go. When a
// custom writer is injected (tests), use it; otherwise write to stdout.
func (p Project) outputWriter() io.Writer {
	if _, isTab := p.writer.(*tabwriter.Writer); isTab {
		return os.Stdout
	}
	return p.writer
}

func (p Project) header() []string {
	return []string{
		"KEY",
		"NAME",
		"TYPE",
		"LEAD",
	}
}

func (p Project) printHeader() {
	headers := p.header()
	end := len(headers) - 1
	for i, h := range headers {
		_, _ = fmt.Fprintf(p.writer, "%s", h)
		if i != end {
			_, _ = fmt.Fprintf(p.writer, "\t")
		}
	}
	_, _ = fmt.Fprintln(p.writer)
}
