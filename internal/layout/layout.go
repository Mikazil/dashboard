package layout

import (
	"github.com/charmbracelet/lipgloss"
	"dashboard/internal/theme"
)

type Cell struct {
	Title string
	View  func(width, height int) string
}

type Grid struct {
	cells [][]Cell
}

func New(cells [][]Cell) *Grid {
	return &Grid{cells: cells}
}

func (g *Grid) View(width, height int) string {
	if len(g.cells) == 0 {
		return ""
	}

	rows := len(g.cells)
	cols := len(g.cells[0])

	rowHeight := height / rows
	colWidth := width / cols

	var rowStrs []string
	for rowIdx, row := range g.cells {
		var colStrs []string
		for colIdx, cell := range row {
			cw := colWidth
			if colIdx == cols-1 {
				cw = width - colWidth*(cols-1)
			}
			ch := rowHeight
			if rowIdx == rows-1 {
				ch = height - rowHeight*(rows-1)
			}

			content := cell.View(cw-2, ch-2)

			box := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(theme.Dim).
				Background(theme.Bg).
				Foreground(theme.Primary).
				Width(cw - 2).
				Render(content)

			colStrs = append(colStrs, box)
		}
		rowStrs = append(rowStrs, lipgloss.JoinHorizontal(lipgloss.Top, colStrs...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rowStrs...)
}
