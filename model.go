package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

type (
	Pos struct {
		X int
		Y int
	}

	XO string

	model struct {
		nextPiece XO
		cursor    Pos
		grid      [N][N]XO
		open      [N][N]bool
	}
)

const (
	// Board Size
	N = 3

	// Define game pieces
	Empty XO = "   "
	X     XO = " X "
	O     XO = " O "
)

var opponent = map[XO]XO{
	X: O,
	O: X,
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	var grid [N][N]XO
	for i := range grid {
		for j := range grid[i] {
			grid[i][j] = Empty
		}
	}

	return model{
		nextPiece: Empty,
		cursor:    Pos{X: -1, Y: -1},
		grid:      grid,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:

		gameOn := Game != nil && !Game.over()
		switch msg.String() {
		case "s":
			if Game == nil {
				m.setGame()
			}
		case "n":
			if Game != nil && Game.over() {
				m = initialModel()
				Game = nil
			}
		case "q":
			return m, tea.Quit
		case "d":
			if gameOn && m.diceAvailable() {
				Game.rollDice()
			}

		case "p":
			if gameOn && Game.dice != nil {
				m.handleDiceResult()
			}

		case "up", "k":
			if gameOn && m.cursor.Y > 0 {
				// if the above square is open, select it.
				if m.open[m.cursor.Y-1][m.cursor.X] {
					m.cursor.Y--
				} else {
					// else find the next open square in the upper matrix
				up_outer_loop:
					for y := m.cursor.Y - 1; y >= 0; y-- {
						for x := range N {
							if m.open[y][x] {
								m.cursor.X = x
								m.cursor.Y = y
								break up_outer_loop
							}
						}
					}
				}
			}
		case "down", "j":
			if gameOn && m.cursor.Y < N-1 {
				if m.open[m.cursor.Y+1][m.cursor.X] {
					m.cursor.Y++
				} else {
				down_outer_loop:
					for y := m.cursor.Y + 1; y < N; y++ {
						for x := range N {
							if m.open[y][x] {
								m.cursor.X = x
								m.cursor.Y = y
								break down_outer_loop
							}
						}
					}
				}
			}

		case "left", "h":
			// skip the used blocks
			if gameOn && m.cursor.X > 0 {
				for x := m.cursor.X - 1; x >= 0; x-- {
					if m.open[m.cursor.Y][x] {
						m.cursor.X = x
						break
					}
				}
			}

		case "right", "l":
			if gameOn && m.cursor.X < N-1 {
				for x := m.cursor.X + 1; x < N; x++ {
					if m.open[m.cursor.Y][x] {
						m.cursor.X = x
						break
					}
				}
			}

		case "space":
			if gameOn && m.open[m.cursor.Y][m.cursor.X] {
				m.grid[m.cursor.Y][m.cursor.X] = m.nextPiece
				m.checkState()
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s := titleUI

	if Game != nil && Game.dice != nil {
		s += m.dice()
	} else {
		s += m.board()
	}

	if Game == nil {
		s += menuUI
	} else if Game.isWon() {
		s += fmt.Sprintf(winUI, string(Game.state))
	} else if Game.state == tie {
		s += tieUI
	} else {
		if Game.dice != nil {
			if Game.winner == XO(Game.state) {
				s += fmt.Sprintf(diceUI, Game.winner, fmt.Sprintf(flipPiece, Game.state))
			} else {
				s += fmt.Sprintf(diceUI, Game.winner, fmt.Sprintf(takeTurn, Game.state))
			}
		} else {
			if m.diceAvailable() {
				s += fmt.Sprintf(turnUI, string(Game.state), diceAvailableUI)
			} else {
				s += fmt.Sprintf(turnUI, string(Game.state), ".")
			}
		}
	}

	return tea.NewView(s)
}
