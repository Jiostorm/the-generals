package network

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"the-generals/board"
	"the-generals/player/soldier"
)

type Handler struct{}

func (Handler) Opponent(player string) string {
	if player == "1" {
		return "2"
	}
	return "1"
}

func (Handler) LoadBoard(b *board.Board, officials map[string][]soldier.Soldier, board_path string) {
	board_file, err := os.ReadFile(board_path)
	if err != nil {
		log.Fatal("File not found!")
	}

	filtered_board := strings.Split(string(board_file), "|")

	id := officials["FL"][0].PlayerID

	i, t := 0, 4
	if id == "1" {
		i, t = 4, board.BOARD_DIMENSION
	}
	for ; i < t; i++ {
		row := strings.Split(strings.TrimSpace(filtered_board[i]), " ")

		for j := range row {
			if row[j] == "--" {
				b.Field[i][j] = &soldier.Soldier{ID: "--"}
			} else {
				b.Field[i][j] = &officials[row[j]][0]
				if len(officials[row[j]][1:]) > 1 {
					officials[row[j]] = officials[row[j]][1:]
				}
			}
		}
	}
}

func (Handler) PrintBoard(board *board.Board, player_id string) {
	fmt.Printf("\n[PLAYER %s]\n", player_id)
	fmt.Println("     1   2   3   4   5   6   7   8")
	fmt.Println("  __________________________________\n  |                                |")

	for i := range board.Field {
		if i == 4 {
			fmt.Println("  | ++++++++++++++++++++++++++++++ |")
		}

		fmt.Print(string(byte('A'+i)), " | ")
		for j := range board.Field[i] {
			if board.Field[i][j].PlayerID == player_id || board.Field[i][j].ID == "--" {
				fmt.Printf("%s  ", board.Field[i][j].ID)
			} else {
				fmt.Print("##  ")
			}
		}
		fmt.Println("\b|", string(byte('A'+i)))
	}
	fmt.Println("  |________________________________|")
	fmt.Println("\n    1   2   3   4   5   6   7   8")
}

func (Handler) LoadOfficials(player_id string, officials_path string) map[string][]soldier.Soldier {
	officials_file, err := os.ReadFile(officials_path)
	if err != nil {
		log.Fatal("File not found!")
	}

	officials := strings.Split(string(officials_file), "|")

	deployed := make(map[string][]soldier.Soldier)
	for i := range officials {
		officials[i] = strings.TrimSpace(officials[i])

		deployee := strings.Split(officials[i], ",")
		rank := deployee[0]
		id := deployee[1]
		count, _ := strconv.Atoi(deployee[2])
		power, _ := strconv.Atoi(deployee[3])
		spyable := false
		if power > 3 || power == 1 {
			spyable = true
		}

		for i := 0; i < count; i++ {
			deployed[id] = append(deployed[id], soldier.Soldier{
				PlayerID:  player_id,
				Rank:      rank,
				ID:        id,
				Count:     count,
				Power:     power,
				IsSpyable: spyable,
			})
		}
	}
	return deployed
}
