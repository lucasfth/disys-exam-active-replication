# disys-exam-active-replication
Made for the 3rd semester subject DISYS

## How to run

### Start servers

First the servers has to be started.
First write command:
```bash
go run server/server.go
```
If it is the first server then write `1` after running the above command. If it is the second then write `2` and if it is the third write `3` and hit enter.
These steps has to be done for all three servers.

### Start client

Then the clients can be started.
Write the command:
```bash
go run client/client.go
```
Then you have to name the client, followed by enter.
This step has to be done for all three clients. Make sure to use a unique name for each client. For the program to create clean log name the client with four characters.

## Crash server
To crash a server you have to write the command `ctrl + c`. This will crash the server entirely.
