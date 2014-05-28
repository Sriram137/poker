package game

import (
	"encoding/json"
	"fmt"
	"github.com/elricL/poker/board"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type CommandMsg struct {
	command string
	name    string
	value   int
}

func sendAll(pokerBoard *board.Board, msg string) {
	var start = pokerBoard.Dealer
	for {
		sendPokerMessage(msg, start.Conn)
		start = start.Next_player
		if start == pokerBoard.Dealer {
			break
		}
	}
}

func sendPokerMessage(msg string, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("EERROEROERO WHile sending message")
	}
}

func HandlePokerMessage(msg []byte, pokerBoard *board.Board, conn *websocket.Conn) {
	var commandMsg map[string]string
	var err = json.Unmarshal(msg, &commandMsg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(commandMsg)
	var name = commandMsg["name"]
	switch commandMsg["command"] {
	case "debug":
		log.Println(pokerBoard)
		pokerBoard.Print()
	case "join":
		if pokerBoard.GameState == "waiting" {
			pokerBoard.AddPlayer(board.Player{nil, false, conn, name, nil, 0, 500})
			if pokerBoard.Length() > 2 {
				pokerBoard.GameState = "canStart"
				sendAll(pokerBoard, "Starting game with "+strconv.Itoa(pokerBoard.Length())+" players")
				gameStart(pokerBoard)
				sendPokerMessage("Your Turn", pokerBoard.Starter.Conn)
			} else {
				sendAll(pokerBoard, name+" Joined")
				sendAll(pokerBoard, "Waiting for "+strconv.Itoa(3-pokerBoard.Length())+" more players")
			}
		} else {
			sendPokerMessage("PLEASE WAIT", conn)
			sendPokerMessage("Chandra is noob", conn)
		}

	case "check":
		log.Println("Check")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			var moneyToCheck = (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - moneyToCheck
			pokerBoard.Pot += moneyToCheck
			pokerBoard.CurrentPlayer.CurrentBet = pokerBoard.CurrentBet
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			bet := strconv.Itoa(pokerBoard.CurrentBet)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" puts in "+strconv.Itoa(moneyToCheck)+" Calls "+bet)
			if pokerBoard.CurrentPlayer == pokerBoard.Starter {
				switch pokerBoard.GameState {
				case "preFlop":
					sendAll(pokerBoard, "Commencing FLOP")
					pokerBoard.GameState = "flop"
					goFlopStuff(pokerBoard)
				case "afterFlop":
					sendAll(pokerBoard, "Commencing TURN")
					pokerBoard.GameState = "afterTurn"
					goTurnStuff(pokerBoard)
				case "afterTurn":
					sendAll(pokerBoard, "Commencing TURN")
					pokerBoard.GameState = "afterRiver"
					goRiverStuff(pokerBoard)
				case "afterRiver":
					sendAll(pokerBoard, "Game Over")
					pokerBoard.GameState = "gameOver"
					findGameWinner(pokerBoard)
				}
			}
		}
	case "raise":
		log.Println("Raise")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			raiseAmount, _ := strconv.Atoi(commandMsg["value"])
			pokerBoard.CurrentPlayer.CurrentBet = raiseAmount
			difference := (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" puts in "+strconv.Itoa(raiseAmount)+" Calls "+strconv.Itoa(pokerBoard.CurrentBet))
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - difference
			pokerBoard.Starter = pokerBoard.CurrentPlayer
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			pokerBoard.Pot += difference
		}
	case "fold":
		log.Println("Fold")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			pokerBoard.CurrentPlayer.Folded = true
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" Folded.")
			nextPlayer := pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			if nextPlayer.FindNextUnfoldedPlayer() == nextPlayer {
				findGameWinner(pokerBoard)
			} else {
				pokerBoard.CurrentPlayer = nextPlayer
			}
		}
	}
	fmt.Println(pokerBoard.GameState)
	fmt.Println()
	// sendAll(msg)
}
