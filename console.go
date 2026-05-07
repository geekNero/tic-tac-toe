package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	titleUI = "Tic-Tac-Toe \n\n"
	menuUI  = "\nPress s to start the game, press q to quit.\n"
	winUI   = "\nWinner %s! Press n to continue.\n"
	tieUI   = "\nGame Tied. Press n to continue\n"
	turnUI  = "\nPlayer%s's turn, press enter or space to lock your move\n"
)

var styler = lipgloss.NewStyle().Reverse(true)

func (m model) cell(x, y int) string {
	value := string(m.grid[y][x])
	if m.cursor.X == x && m.cursor.Y == y {
		return styler.Render(value)
	}
	return value
}

func (m model) board() string {
	s := ""

	for i := range N {

		s += fmt.Sprintf("%s|%s|%s\n", m.cell(0, i), m.cell(1, i), m.cell(2, i))
		if i < 2 {
			s += "-----------\n"
		}
	}

	return s
}
