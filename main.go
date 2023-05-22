package main

import (
	. "github.com/S-A-RB05/StratService/stratserver"
)

const grpcPort = "10000"

func main() {
	InitGRPC(grpcPort)
}
