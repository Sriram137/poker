git clone in GOPATH.
go build && ./poker to run the server.

Install web socket handler using

sudo npm install -g ws


connect to the localhost server using wscat -c ws:localhost:9000/ws

Test using 

{"command" : "join", "name" : "vivek"}
{"command" : "join", "name" : "vivek2"}
{"command" : "join", "name" : "vivek3"}

{"command" : "check", "name" : "vivek"}
{"command" : "check", "name" : "vivek2"}
{"command" : "raise", "name" : "vivek3", "value":"100"}
{"command" : "debug"}
{"command" : "check", "name" : "vivek"}
{"command" : "check", "name" : "vivek2"}

{"command" : "check", "name" : "vivek2"}
{"command" : "check", "name" : "vivek3"}
{"command" : "check", "name" : "vivek"}
