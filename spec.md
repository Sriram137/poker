You'll be prompted with messages "Your Turn" when its your turn

join with 
j <name>
On any player join:
Server returns:
    2 messages.
        1st Message -> <name> joined.
        2nd Message -> Waiting for <3 - number of players> more players.
    if number of players > 2:
    Server returns:
        6 messages.
            1st Message -> Starting game with <number of players> players.
            2nd Message -> Commencing PREFLOP
            3rd Message -> 1st card
            4th Message -> 2nd card
            5th Message -> Your Position <your position>
            6th Message -> Player list in format. ***2 bet:0money :500 -> 3 bet:10money :490 -> 1 bet:20money :480

When someone checks:
Server returns:
    2 messages.
        1st Message -> <name> puts in <amount> to call <current_bet>
        2nd Message -> <name>'s Turn

When someone folds:
Server returns:
    2 messages.
        1st Message -> <name> folded
        2nd Message -> <name>'s Turn

When Flop starts:
Server returns:
    5 messages.
        1st Message -> Commencing FLOP
        2nd Message -> Card
        3rd Message -> Card
        4th Message -> Card
        5th Message -> <name>'s Turn

When Turn starts:
    3 messages.
        1st Message -> Commencing TURN
        2nd Message -> Card
        3rd Message -> <name>'s Turn

When River starts:
    3 messages.
        1st Message -> Commencing RIVER
        2nd Message -> Card
        3rd Message -> <name>'s Turn

When Game Over:
    3 Actions.
        1 message.
            1st Message -> Game Over
        x messages.
            if one winner:
                1 message.
                1st Message -> <name> wins the pot of <amount won>
            if more winners:
                1st Message -> Pot is split between <winner_count> players
                for winner in winners:
                    message -> <name> wins <amount_won>
