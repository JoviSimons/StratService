package stratserver

import (
	"context"
	"fmt"
	"github.com/S-A-RB05/StratService/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

type StratServiceServer struct {
	proto.UnimplementedStratServiceServer
}

func (s StratServiceServer) mustEmbedUnimplementedStratServiceServer() {
}

func (s StratServiceServer) ReturnAll(req *proto.ReturnAllReq, server proto.StratService_ReturnAllServer) error {
	//TODO implement me
	panic("implement me")
}

func (s StratServiceServer) ReturnStrat(ctx context.Context, req *proto.ReturnStratReq) (*proto.ReturnStratRes, error) {
	//TODO implement me
	panic("implement me")
}

func (s StratServiceServer) StoreStrat(ctx context.Context, req *proto.StoreStratReq) (*proto.StoreStratRes, error) {
	//TODO implement me
	panic("implement me")
}

type StratItem struct {
	Name string `bson:""`
}

var db *mongo.Client
var stratdb *mongo.Collection
var mongoCtx context.Context

func InitGRPC() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Starting server on port: 50051")

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// Set options, here we can configure things like TLS support
	var opts []grpc.ServerOption
	// Create new gRPC server with (blank) options
	s := grpc.NewServer(opts...)
	// Create BlogService type
	srv := &StratServiceServer{}
	// Register the service with the server
	proto.RegisterStratServiceServer(s, srv)

	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")

	// non-nil empty context
	mongoCtx = context.Background()

	// Connect takes in a context and options, the connection URI is the only option we pass for now
	// mongodb+srv://stockbrood:admin@stockbrood.sifn3lq.mongodb.net/test
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	// Handle potential errors
	if err != nil {
		log.Fatal(err)
	}

	// Check whether the connection was succesful by pinging the MongoDB server
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	// Bind our collection to our global variable for use in other methods
	stratdb = db.Database("testing").Collection("strategies")

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	// Right way to stop the server using a SHUTDOWN HOOK
	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-c

	// After receiving CTRL+C Properly stop the server
	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}
