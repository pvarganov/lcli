package format

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// TableWriter выводит таблицу в заданный writer.
func TableWriter(w io.Writer, headers []string, rows [][]string) {
	table := tablewriter.NewTable(w)
	table.Header(headers)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}
