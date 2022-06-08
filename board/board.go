package board

import (
	"fmt"
	"math"
	"the-generals/player/soldier"
)

const (
	BOARD_DIMENSION = 8
)

type Board struct {
	Field [][]*soldier.Soldier
}

func (board *Board) Move(moveset []string, player string, opponent string) (string, map[string]string) {
	x := int(byte(moveset[0][0]) - 'A')
	y := int(byte(moveset[0][1]) - '1')

	new_x := int(byte(moveset[1][0]) - 'A')
	new_y := int(byte(moveset[1][1]) - '1')

	if int(math.Abs(float64(x-new_x))+math.Abs(float64(y-new_y))) != 1 {
		return "Invalid Move!", nil
	}
	if board.Field[x][y].PlayerID != player {
		return "Soldier is not Yours!", nil
	}

	if board.Field[new_x][new_y].PlayerID == player {
		return "Spot already occupied by one of your Soldier!", nil
	}

	pvp := board.Field[x][y].Challenge(board.Field[new_x][new_y])

	var results = make(map[string]string)

	switch pvp {
	case 0: // Draw
		board.Field[x][y].Count--
		board.Field[new_x][new_y].Count--

		results[player] = fmt.Sprintf("Your %s is in draw match!", board.Field[x][y].Rank)
		results[opponent] = fmt.Sprintf("Your %s is in draw match!", board.Field[new_x][new_y].Rank)

		board.Field[x][y] = &soldier.Soldier{ID: "--"}
		board.Field[new_x][new_y] = &soldier.Soldier{ID: "--"}
		break
	case 1: // Win
		board.Field[new_x][new_y].Count--

		if board.Field[new_x][new_y].ID != "--" {
			results[player] = fmt.Sprintf("Your %s won!", board.Field[x][y].Rank)
			results[opponent] = fmt.Sprintf("Your %s lost!", board.Field[new_x][new_y].Rank)
		}

		board.Field[new_x][new_y] = board.Field[x][y]
		board.Field[x][y] = &soldier.Soldier{ID: "--"}
		break
	case -1: // Lose
		board.Field[x][y].Count--

		if board.Field[x][y].ID != "--" {
			results[player] = fmt.Sprintf("Your %s lost!", board.Field[x][y].Rank)
			results[opponent] = fmt.Sprintf("Your %s won!", board.Field[new_x][new_y].Rank)
		}
		board.Field[x][y] = &soldier.Soldier{ID: "--"}
		break
	}

	return "G", results
}
