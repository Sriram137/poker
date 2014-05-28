# Devlopment
Clone so that directory structure looks like 
> ~GOPATH/src/github.com/elricL/poker

# Usage 
Need a webscocket client library to connect to the server. 

## Installation

A websocket client library written in node can be installed using

    npm install -g ws

Connect using 

    wscat -c ws:server_address:9000/ws
  
## Test

    j <name> or join <name>
    c or check
    f or fold
    r <amount> or raise <amount>
