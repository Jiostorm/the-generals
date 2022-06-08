package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"the-generals/network"
)

var server network.Server

var handler network.Handler

var flag bool

var turn int

func main() {
	server = network.Server{IP: network.DEFAULT_IP, Port: network.DEFAULT_PORT, Clients: make(map[string]*net.Conn), ConnCount: 0}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.IP, server.Port))
	if err != nil {
		panic(err)
	}
	log.Printf("[SERVER]\tListening at %s:%s\n", server.IP, server.Port)
	fmt.Println("\n==========================================================================")
	defer listener.Close()

	flag = false
	for {
		conn, _ := listener.Accept()
		log.Printf("[SERVER]\tConnection Detected! %s\n", conn.RemoteAddr().String())

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	server.ConnCount++
	io.WriteString(conn, fmt.Sprint(server.ConnCount, "~"))

	if server.ConnCount != 2 {
		log.Printf("[SERVER]\tWaiting for another client...\n")
	} else {
		log.Printf("[SERVER]\tThe Game of Generals will start!\n\n")
	}

	server.Clients[strconv.Itoa(server.ConnCount)] = &conn

	for {
		if server.ConnCount == 2 {
			break
		}
	}

	defer conn.Close()

	player := ""

	if !flag {
		io.WriteString(*server.Clients["1"], fmt.Sprint("[turn] player 1~"))
		flag = true
		turn = 0
	}

	for flag && server.ConnCount == 2 {
		signal, err := bufio.NewReader(conn).ReadString('~')
		signal_args := strings.Split(strings.TrimSpace(signal[:len(signal)-1]), " ")
		if err != nil {
			panic(err)
		}

		header := signal_args[0]
		player = signal_args[2]

		if header == "[move]" {
			turn++
			moveset := signal_args[3:]
			flag_count := signal_args[5:]
			log.Printf("[PLAYER %s]\tTurn %02d :: Moved a Soldier (%s -> %s)", player, turn, moveset[0], moveset[1])
			log.Printf("[SERVER]\tForwading reply to [PLAYER %s]\n", handler.Opponent(player))

			// Replying to another Player
			io.WriteString(*server.Clients[handler.Opponent(player)], fmt.Sprintf("[move] player %s %s %s~", player, moveset[0], moveset[1]))

			if flag_count[0] != flag_count[1] {
				winner := ""
				if flag_count[0] == "0" {
					winner = handler.Opponent(player)
				} else if flag_count[1] == "0" {
					winner = player
				}
				io.WriteString(*server.Clients[player], fmt.Sprintf("[set] player %s~", winner))
				io.WriteString(*server.Clients[handler.Opponent(player)], fmt.Sprintf("[set] player %s~", winner))

				log.Printf("[SERVER]\t[Player %s] won the match!", winner)
				fmt.Print("==========================================================================\n\n")
				log.Println("[SERVER]\tAccepting new clients!")
				fmt.Print("==========================================================================\n\n")
				break
			}
			continue
		}
		// Giving turn to another Player
		io.WriteString(*server.Clients[player], fmt.Sprintf("[turn] player %s~", player))
	}
	server.ConnCount--
	flag = false
}
