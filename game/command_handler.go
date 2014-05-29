package game

import (
	"github.com/elricL/poker/board"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"strings"
)

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

func getRequestingPlayer(pokerBoard *board.Board, conn *websocket.Conn) *board.Player {
	var player = pokerBoard.Dealer
	for {
		if player.Conn == conn {
			log.Println("identified " + player.Name)
			return player
		}
		player = player.Next_player
		if player == pokerBoard.Dealer {
			return nil
		}
	}
}

func sendPokerMessage(msg string, conn *websocket.Conn) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		log.Println("EERROEROERO WHile sending message")
	}
}

func getCommand(stringMsg string) (string, string) {
	splits := strings.Split(stringMsg, " ")
	command := ""
	value := ""
	switch len(splits) {
	case 1:
		command = splits[0]
	case 2:
		command = splits[0]
		value = splits[1]
	}
	switch command {
	case "j":
		command = "join"
	case "f":
		command = "fold"
	case "c":
		command = "check"
	case "r":
		command = "raise"
	}
	return command, value
}

func HandlePokerMessage(msg []byte, pokerBoard *board.Board, conn *websocket.Conn) {
	var stringMsg = string(msg)
	command, command_value := getCommand(stringMsg)
	switch command {
	case "debug":
		log.Println(pokerBoard)
		pokerBoard.Print()
	case "join":
		if pokerBoard.Length() == 0 {
			pokerBoard.BoardCards = []string{"__", "__", "__", "__", "__"}
		}
		if pokerBoard.Length() > 0 && getRequestingPlayer(pokerBoard, conn) != nil {
			sendPokerMessage("We know you are over enthusiastic about Poker.But only one instance of you can join a table!!", conn)
			return
		}
		if pokerBoard.GameState == "waiting" {
			pokerBoard.AddPlayer(board.Player{nil, false, conn, command_value, nil, 0, 500})
			if pokerBoard.Length() > 2 {
				pokerBoard.GameState = "canStart"
				sendAll(pokerBoard, "Starting game with "+strconv.Itoa(pokerBoard.Length())+" players")
				gameStart(pokerBoard)
				sendPokerMessage("Your Turn", pokerBoard.Starter.Conn)
			} else {
				sendAll(pokerBoard, command_value+" Joined")
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
		log.Println(pokerBoard.CurrentPlayer.Next_player)
		log.Println(pokerBoard.CurrentPlayer.Next_player.Next_player)
		if pokerBoard.CurrentPlayer.Conn == conn {
			var moneyToCheck = (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - moneyToCheck
			pokerBoard.Pot += moneyToCheck
			pokerBoard.CurrentPlayer.CurrentBet = pokerBoard.CurrentBet
			bet := strconv.Itoa(pokerBoard.CurrentBet)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" puts in "+strconv.Itoa(moneyToCheck)+" calls "+bet)
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
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
					sendAll(pokerBoard, "Commencing RIVER")
					pokerBoard.GameState = "afterRiver"
					goRiverStuff(pokerBoard)
				case "afterRiver":
					sendAll(pokerBoard, "Game Over")
					pokerBoard.GameState = "gameOver"
					findGameWinner(pokerBoard)
				}
			}
			sendPokerMessage("Your Turn", pokerBoard.CurrentPlayer.Conn)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+"'s Turn")
		} else {
			sendPokerMessage("Out of Turn", conn)
		}
		log.Println(pokerBoard.CurrentPlayer)
		log.Println(pokerBoard.CurrentPlayer.Next_player)
		log.Println(pokerBoard.CurrentPlayer.Next_player.Next_player)
	case "raise":
		log.Println("Raise")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		log.Println(pokerBoard.CurrentPlayer.Next_player)
		log.Println(pokerBoard.CurrentPlayer.Next_player.Next_player)
		if pokerBoard.CurrentPlayer.Conn == conn {
			raiseAmount, _ := strconv.Atoi(command_value)
			difference := (raiseAmount - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.CurrentBet = raiseAmount
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" puts in "+strconv.Itoa(difference)+" raises to "+strconv.Itoa(raiseAmount))
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - difference
			pokerBoard.Starter = pokerBoard.CurrentPlayer
			pokerBoard.CurrentBet = raiseAmount
			pokerBoard.Pot += difference
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			sendPokerMessage("Your Turn", pokerBoard.CurrentPlayer.Conn)
		} else {
			sendPokerMessage("Out of turn", conn)
		}
		log.Println(pokerBoard.CurrentPlayer)
		log.Println(pokerBoard.CurrentPlayer.Next_player)
		log.Println(pokerBoard.CurrentPlayer.Next_player.Next_player)
	case "fold":
		log.Println("Fold")
		log.Println(pokerBoard)
		log.Println(pokerBoard.CurrentPlayer)
		if pokerBoard.CurrentPlayer.Conn == conn {
			pokerBoard.CurrentPlayer.Folded = true
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" Folded.")
			if pokerBoard.CurrentPlayer == pokerBoard.Starter {
				pokerBoard.Starter = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			}
			nextPlayer := pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			if nextPlayer.FindNextUnfoldedPlayer() == nextPlayer {
				log.Println("Winner Winner Chicken Dinner")
				findGameWinner(pokerBoard)
			} else {
				pokerBoard.CurrentPlayer = nextPlayer
				sendPokerMessage("Your Turn", pokerBoard.CurrentPlayer.Conn)
			}
		} else {
			sendPokerMessage("Out of turn", conn)
		}

	case "me":
		log.Println("me")
		var player = getRequestingPlayer(pokerBoard, conn)
		sendPokerMessage(player.PlayerInfo(), player.Conn)

	case "board":
		log.Println("board")
		sendPokerMessage(pokerBoard.PrintCards(), getRequestingPlayer(pokerBoard, conn).Conn)

	case "who":
		log.Println("who")
		var player = getRequestingPlayer(pokerBoard, conn)
		if pokerBoard.GameState == "waiting" {
			sendPokerMessage("Patience My friend. Let the Game begin", player.Conn)
		} else {
			if player == pokerBoard.CurrentPlayer {
				sendPokerMessage("Confused eh!!. Its YOUR turn", player.Conn)
			} else {
				sendPokerMessage(pokerBoard.CurrentPlayer.Name+"'s turn", player.Conn)
			}
		}
	case "players":
		//TODO EdgeCases need to be checked ,Sriram
		log.Println("players")
		var player = getRequestingPlayer(pokerBoard, conn)
		if pokerBoard.GameState == "waiting" {
			sendPokerMessage("Patience My friend. Let the Game begin", player.Conn)
			return
		}
		var currPlayer = pokerBoard.CurrentPlayer
		sendPokerMessage("Current Player is "+currPlayer.Name+" with bet "+strconv.Itoa(currPlayer.CurrentBet), player.Conn)
		for {
			currPlayer = currPlayer.FindNextUnfoldedPlayer()
			if currPlayer == pokerBoard.CurrentPlayer {
				break
			}
			sendPokerMessage("Next Player is "+currPlayer.Name+" with bet "+strconv.Itoa(currPlayer.CurrentBet), player.Conn)
		}

	}

	log.Println(pokerBoard.GameState)
	log.Println()
	// sendAll(msg)
}
