package gameModel

func checkWinner(row []string) bool {
	w := row[0]
	for _, v := range row {
		if w == "" || w != v {
			return false
		}
	}
	return true
}

func getWinner(board [][]string) string {
	l := len(board)
	row := make([]string, l)
	for i := 0; i < l; i++ {
		if checkWinner(board[i]) {
			return board[i][0]
		}
		for j := 0; j < l; j++ {
			row[j] = board[j][i]
		}
		if checkWinner(row) {
			return row[0]
		}
	}
	for i := 0; i < l; i++ {
		row[i] = board[i][i]
	}
	if checkWinner(row) {
		return row[0]
	}
	for i := 0; i < l; i++ {
		row[i] = board[2-i][i]
	}
	if checkWinner(row) {
		return row[0]
	}
	return ""
}
