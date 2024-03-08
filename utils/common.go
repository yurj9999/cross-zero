package utils

import (
	"fmt"
	"net/http"
)

func CalcWinner(player string, fields [][]string) string {
	firstHorizontalWin := fields[0][0] == player && fields[0][1] == player && fields[0][2] == player
	secondHorizontalWin := fields[1][0] == player && fields[1][1] == player && fields[1][2] == player
	thirdHorizontalWin := fields[2][0] == player && fields[2][1] == player && fields[2][2] == player

	firstVerticallWin := fields[0][0] == player && fields[1][0] == player && fields[2][0] == player
	secondVerticalWin := fields[1][0] == player && fields[1][1] == player && fields[1][2] == player
	thirdVerticalWin := fields[2][0] == player && fields[2][1] == player && fields[2][2] == player

	firstDiagonalWin := fields[0][0] == player && fields[1][1] == player && fields[2][2] == player
	secondDiagonalWin := fields[0][2] == player && fields[1][1] == player && fields[2][0] == player

	if firstHorizontalWin || secondHorizontalWin || thirdHorizontalWin || firstVerticallWin || secondVerticalWin || thirdVerticalWin || firstDiagonalWin || secondDiagonalWin {
		return player
	}

	return ""
}

func GetFieldStatus(fields [][]string, w http.ResponseWriter) {
	for _, data := range fields {
		for indexItem, item := range data {
			if indexItem == 1 {
				fmt.Fprintf(w, "| %s |", item)
			} else if indexItem == LAST_POSITION {
				fmt.Fprintf(w, " %s\n", item)
			} else {
				fmt.Fprintf(w, "%s ", item)
			}
		}

		fmt.Fprintln(w, "----------")
	}
}
