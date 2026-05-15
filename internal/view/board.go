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

// BoardOption is a functional option to wrap board properties.
type BoardOption func(*Board)

// Board is a board view.
type Board struct {
	data      []*jira.Board
	writer    io.Writer
	buf       *bytes.Buffer
	plain     bool
	noHeaders bool
	csv       bool
	delimiter string
}

// NewBoard initializes a board.
func NewBoard(data []*jira.Board, opts ...BoardOption) *Board {
	b := Board{
		data: data,
		buf:  new(bytes.Buffer),
	}
	b.writer = tabwriter.NewWriter(b.buf, 0, tabWidth, 1, '\t', 0)

	for _, opt := range opts {
		opt(&b)
	}
	return &b
}

// WithBoardWriter sets a writer for the board.
func WithBoardWriter(w io.Writer) BoardOption {
	return func(b *Board) {
		b.writer = w
	}
}

// WithBoardPlain toggles plain output (no interactive pager).
func WithBoardPlain(v bool) BoardOption {
	return func(b *Board) {
		b.plain = v
	}
}

// WithBoardNoHeaders skips the header row.
func WithBoardNoHeaders(v bool) BoardOption {
	return func(b *Board) {
		b.noHeaders = v
	}
}

// WithBoardCSV emits CSV output.
func WithBoardCSV(v bool) BoardOption {
	return func(b *Board) {
		b.csv = v
	}
}

// WithBoardDelimiter sets the column delimiter in plain mode. Default is "\t".
func WithBoardDelimiter(d string) BoardOption {
	return func(b *Board) {
		b.delimiter = d
	}
}

// Render renders the board view.
func (b Board) Render() error {
	data := b.tableData()

	if b.csv {
		return renderCSV(b.outputWriter(), data)
	}

	delim := b.delimiter
	if delim == "" {
		delim = "\t"
	}
	if b.plain && delim != "\t" {
		return renderPlain(b.outputWriter(), data, delim)
	}

	for _, row := range data {
		for i, cell := range row {
			_, _ = fmt.Fprint(b.writer, cell)
			if i != len(row)-1 {
				_, _ = fmt.Fprint(b.writer, "\t")
			}
		}
		_, _ = fmt.Fprintln(b.writer)
	}
	if tw, ok := b.writer.(*tabwriter.Writer); ok {
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	if _, ok := b.writer.(*tabwriter.Writer); !ok {
		return nil
	}
	if b.plain {
		_, err := fmt.Fprint(os.Stdout, b.buf.String())
		return err
	}
	return tui.PagerOut(b.buf.String())
}

func (b Board) tableData() tui.TableData {
	var data tui.TableData
	if !b.noHeaders {
		data = append(data, b.header())
	}
	for _, d := range b.data {
		data = append(data, []string{fmt.Sprintf("%d", d.ID), prepareTitle(d.Name), d.Type})
	}
	return data
}

func (b Board) outputWriter() io.Writer {
	if _, isTab := b.writer.(*tabwriter.Writer); isTab {
		return os.Stdout
	}
	return b.writer
}

func (b Board) header() []string {
	return []string{
		"ID",
		"NAME",
		"TYPE",
	}
}

func (b Board) printHeader() {
	n := len(b.header())
	for i, h := range b.header() {
		_, _ = fmt.Fprintf(b.writer, "%s", h)
		if i != n-1 {
			_, _ = fmt.Fprintf(b.writer, "\t")
		}
	}
	_, _ = fmt.Fprintln(b.writer, "")
}
