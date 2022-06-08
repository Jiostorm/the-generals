package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"the-generals/board"
	"the-generals/network"
	"the-generals/player/soldier"
)

const (
	OFFICIALS_FILE = "officials.txt"
	BOARD_FILE     = "board.txt"
)

var client network.Client

var my_officers, opp_officers map[string][]soldier.Soldier

var game_board board.Board

var handler network.Handler

var player string

func start() { // Loading of Player configurations
	temp := make([][]*soldier.Soldier, board.BOARD_DIMENSION)
	for i := range temp {
		temp[i] = make([]*soldier.Soldier, board.BOARD_DIMENSION)
	}

	my_officers = handler.LoadOfficials(player, OFFICIALS_FILE)
	opp_officers = handler.LoadOfficials(handler.Opponent(player), OFFICIALS_FILE)

	game_board = board.Board{Field: temp}

	handler.LoadBoard(&game_board, my_officers, BOARD_FILE)
	handler.LoadBoard(&game_board, opp_officers, BOARD_FILE)
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Insufficient Arguments!")
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	player, _ = bufio.NewReader(conn).ReadString('~')
	player = player[:len(player)-1]
	client = network.Client{ID: player, Conn: conn}

	start() // Load the board

	fmt.Println("\n============================================")
	fmt.Println("Your Starting Board:")
	handler.PrintBoard(&game_board, player)
	fmt.Println("============================================")

	for {
		signal, err := bufio.NewReader(conn).ReadString('~')
		signal_args := strings.Split(strings.TrimSpace(signal)[:len(signal)-1], " ")
		signal_len := len(signal_args)
		if err != nil {
			panic(err)
		}

		header := signal_args[0]
		player = signal_args[2]

		switch header {
		case "[move]":
			if signal_len != 5 {
				log.Fatal("Insufficient Moveset!")
				continue
			}
			moveset := signal_args[3:]
			fmt.Println("\n============================================")
			fmt.Printf("Opposing Player Moved a Soldier (%s -> %s)\n", moveset[0], moveset[1])

			_, results := game_board.Move(moveset, player, client.ID)

			fmt.Println("--------------------------------------------")
			fmt.Println(results[client.ID])
			handler.PrintBoard(&game_board, client.ID)
			fmt.Println("============================================")

			io.WriteString(client.Conn, fmt.Sprintf("[turn] player %s~", client.ID))
			break
		case "[turn]":
			if player != client.ID {
				continue
			}
			fmt.Println("Your Turn!\n--------------------------------------------")

			move := ""

			for {
				fmt.Print("C: Current Position\nT: Target Position\nEnter Your Move(C T): ")
				move, _ = bufio.NewReader(os.Stdin).ReadString('\n')
				move = strings.TrimSpace(move)

				fmt.Println("--------------------------------------------")

				valid, results := game_board.Move(strings.Split(move, " "), client.ID, handler.Opponent(client.ID))
				if valid == "G" {
					if results != nil {
						fmt.Println(results[client.ID])
					}
					break
				}
				fmt.Printf("[%s]\n", valid)
			}
			flag_count := []int{my_officers["FL"][0].Count, opp_officers["FL"][0].Count}

			handler.PrintBoard(&game_board, client.ID)

			io.WriteString(client.Conn, fmt.Sprintf("[move] player %s %s %d %d~", client.ID, move, flag_count[0], flag_count[1]))
			fmt.Println("============================================")
			fmt.Println("Waiting...")
			break
		case "[set]":
			fmt.Println("\nGame is Done!")
			switch player {
			case client.ID:
				fmt.Println("-------------\nYOU WIN\n-------------")
				break
			case handler.Opponent(client.ID):
				fmt.Println("-------------\nYOU LOSE\n-------------")
			}
			return
		}
	}
}
