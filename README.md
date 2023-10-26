# duplex-stress
Full-duplex stress scenario that mocks a multiplexed protocol 

## Building

1. Install Golang 1.21+ and Powershell
2. .\build.ps1

Produces a Linux and Win executable 

## Running 

In one terminal

`./duplex-stress server 0.0.0.0 9191`

In another

`./duplex-stress client 0.0.0.0 9191`

If the connection is working it will keep repeating until either side is ctrl-c/z'd. Otherwise it will hang.

Note: Server only answers to a single client
