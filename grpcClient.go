package main

import (
	"context"
	"fmt"
	"github.com/S-A-RB05/StratService/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	c := proto.NewStratServiceClient(conn)

	fmt.Println(ReturnAll(c))
}

func ReturnAll(c proto.StratServiceClient) (response proto.StratService_ReturnAllClient) {
	response, err := c.ReturnAll(context.Background(), &proto.ReturnAllReq{})
	if err != nil {
		log.Fatalf("Failed to call MyRPCMethod: %v", err)
	}
	return response
}

func StoreStrat()
