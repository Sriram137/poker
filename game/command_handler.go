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

func getPlayerPositionsInPrintableFormat(pokerBoard *board.Board) string {
	var currPlayer = pokerBoard.CurrentPlayer
	var playerPosString string = ""
	for {
		playerPosString = playerPosString + currPlayer.Name + " "
		if currPlayer.Folded == true {
			playerPosString = playerPosString + "(f) "
		}
		playerPosString = playerPosString + "bet:" + strconv.Itoa(currPlayer.CurrentBet) + "money :" + strconv.Itoa(currPlayer.Money) + " "
		if currPlayer == pokerBoard.Starter {
			playerPosString = "***" + playerPosString
		}
		if currPlayer.Next_player != pokerBoard.CurrentPlayer {
			playerPosString = playerPosString + "-> "
		}
		currPlayer = currPlayer.Next_player
		if currPlayer == pokerBoard.CurrentPlayer {
			return playerPosString
		}
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
	case "b":
		command = "board"
	case "p":
		command = "pot"
	case "ls":
		command = "players"
	case "w":
		command = "who"
	}
	return command, value
}

func finishGame(pokerBoard *board.Board) {
	winners, amount := findGameWinner(pokerBoard)
	if len(winners) == 1 {
		sendAll(pokerBoard, strings.Join([]string{winners[0].Name, "wins the pot of", strconv.Itoa(amount)}, " "))
	} else {
		sendAll(pokerBoard, strings.Join([]string{"Pot is split between", strconv.Itoa(len(winners)), "players"}, " "))
		for _, winner := range winners {
			sendAll(pokerBoard, strings.Join([]string{winner.Name, "wins", strconv.Itoa(amount)}, " "))
		}
	}
	resetGame(pokerBoard)
}

func goNextState(pokerBoard *board.Board) {
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
			finishGame(pokerBoard)
		}
	}
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
		if command_value == "" {
			sendPokerMessage("Enter a name", conn)
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
			pokerBoard.AddPlayer(board.Player{nil, true, conn, command_value, nil, 0, 500})
			sendPokerMessage("Please wait till the current game is over", conn)
			sendPokerMessage("Chandra is noob", conn)
		}
	case "check":
		log.Println("Check")
		if pokerBoard.CurrentPlayer.Conn == conn {
			var moneyToCheck = (pokerBoard.CurrentBet - pokerBoard.CurrentPlayer.CurrentBet)
			pokerBoard.CurrentPlayer.Money = pokerBoard.CurrentPlayer.Money - moneyToCheck
			pokerBoard.Pot += moneyToCheck
			pokerBoard.CurrentPlayer.CurrentBet = pokerBoard.CurrentBet
			bet := strconv.Itoa(pokerBoard.CurrentBet)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" puts in "+strconv.Itoa(moneyToCheck)+" calls "+bet)
			pokerBoard.CurrentPlayer = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			goNextState(pokerBoard)
			sendPokerMessage("Your Turn", pokerBoard.CurrentPlayer.Conn)
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+"'s Turn")
		} else {
			sendPokerMessage("Out of Turn", conn)
		}
	case "raise":
		log.Println("Raise")
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
	case "fold":
		log.Println("Fold")
		if pokerBoard.CurrentPlayer.Conn == conn {
			pokerBoard.CurrentPlayer.Folded = true
			sendAll(pokerBoard, pokerBoard.CurrentPlayer.Name+" Folded.")
			if pokerBoard.CurrentPlayer == pokerBoard.Starter {
				pokerBoard.Starter = pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			}
			nextPlayer := pokerBoard.CurrentPlayer.FindNextUnfoldedPlayer()
			if nextPlayer.FindNextUnfoldedPlayer() == nextPlayer {
				log.Println("Winner Winner Chicken Dinner")
				finishGame(pokerBoard)
			} else {
				pokerBoard.CurrentPlayer = nextPlayer
				goNextState(pokerBoard)
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
		sendPokerMessage(getPlayerPositionsInPrintableFormat(pokerBoard), conn)
	case "pot":
		log.Println("pot")
		sendPokerMessage("The current pot is "+strconv.Itoa(pokerBoard.Pot), conn)
	default:
		sendPokerMessage("Go to https://github.com/elricL/poker/blob/master/README.md for usage instructions", conn)
	}

	log.Println(pokerBoard.GameState)
	log.Println()
	// sendAll(msg)
}
