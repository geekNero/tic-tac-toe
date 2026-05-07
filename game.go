package main

type (
	state XO

	game struct {
		state
	}
)

const (
	playingO = state(O)
	playingX = state(X)

	// win and playing won't conflict as X and O are padded by spaces
	winO = state("O")
	winX = state("X")
	tie  = state("tie")
)

var Game *game

func (m *model) setGame() {
	Game = &game{
		state: playingO,
	}

	*m = initialModel()
	m.nextPiece = O
	m.cursor.X = 0
	m.cursor.Y = 0

	// open the grid
	m.open = [N][N]bool{}

	for i := range N {
		for j := range N {
			m.open[i][j] = true
		}
	}
}

func (m *model) checkState() {
	m.open[m.cursor.Y][m.cursor.X] = false

	switch Game.state {
	case playingO:
		if m.checkWinCondition(O) {
			Game.state = winO
			return
		}

		Game.state = playingX
		m.nextPiece = X

	case playingX:
		if m.checkWinCondition(X) {
			Game.state = winX
			return
		}

		Game.state = playingO
		m.nextPiece = O
	}

	if m.gridExhausted() {
		Game.state = tie
	}
}

func (m *model) checkWinCondition(player XO) bool {
	x := m.cursor.X

	// check y - axis
	wonYAxis := true
	for y := range N {
		if m.grid[y][x] != player {
			wonYAxis = false
			break
		}
	}

	// check x - axis
	wonXAxis := true
	y := m.cursor.Y
	for x := range N {
		if m.grid[y][x] != player {
			wonXAxis = false
			break
		}
	}

	// check upper diagonal
	wonDiagonal1 := true
	for y, x := 0, 0; y < N && x < N; {
		if m.grid[y][x] != player {
			wonDiagonal1 = false
			break
		}
		y++
		x++
	}

	// check lower diagonal
	wonDiagonal2 := true
	for y, x := N-1, 0; y >= 0 && x < N; {
		if m.grid[y][x] != player {
			wonDiagonal2 = false
			break
		}
		y--
		x++
	}

	return wonYAxis || wonXAxis || wonDiagonal1 || wonDiagonal2
}

func (g game) isWon() bool {
	return g.state == winO || g.state == winX
}

func (g game) over() bool {
	return g.state == winO || g.state == winX || g.state == tie
}

func (m model) gridExhausted() bool {
	for y := range N {
		for x := range N {
			if m.open[y][x] {
				return false
			}
		}
	}

	return true
}
