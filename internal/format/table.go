package format

import (
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

// Table выводит данные в виде таблицы.
func Table(headers []string, rows [][]string) {
	TableWriter(os.Stdout, headers, rows)
}

// TableWriter выводит таблицу в заданный writer.
func TableWriter(w io.Writer, headers []string, rows [][]string) {
	table := tablewriter.NewTable(w)
	table.Header(headers)
	for _, row := range rows {
		table.Append(row)
	}
	table.Render()
}
