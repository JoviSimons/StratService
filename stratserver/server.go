package stratserver

import (
	"context"
	"fmt"
	"github.com/S-A-RB05/StratService/proto"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

type StratServiceServer struct {
	proto.UnimplementedStratServiceServer
}

func (s StratServiceServer) ReturnAll(req *proto.ReturnAllReq, server proto.StratService_ReturnAllServer) error {
	data := &StratItem{}
	cursor, err := stratdb.Find(context.TODO(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&data)
		if err != nil {
			log.Printf("error decoding data: %v", err)
			continue
			//return status.Errorf(codes.Internal, fmt.Sprintf("Could not decode data: %v", err))
		}

		server.Send(&proto.ReturnAllRes{
			Strategy: &proto.Strategy{
				Name:    data.Name,
				Mq:      data.Mq,
				Ex:      data.Ex,
				Created: timestamppb.New(data.Created),
			},
		})
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown cursor error: %v", err))
	}
	return nil
}

func (s StratServiceServer) ReturnStrat(ctx context.Context, req *proto.ReturnStratReq) (*proto.ReturnStratRes, error) {
	fmt.Println("searching for strat with name " + req.Name)
	result := stratdb.FindOne(ctx, bson.M{"_name": req.GetName()})
	data := StratItem{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find Strategy with name %s: %v", req.GetName(), err))
	}
	response := &proto.ReturnStratRes{
		Strategy: &proto.Strategy{
			Name:    data.Name,
			Mq:      data.Mq,
			Ex:      data.Ex,
			Created: timestamppb.New(data.Created),
		},
	}
	return response, nil
}

func (s StratServiceServer) StoreStrat(_ context.Context, req *proto.StoreStratReq) (*proto.StoreStratRes, error) {
	strategy := req.GetStrategy()
	data := StratItem{
		Name:    strategy.GetName(),
		Mq:      strategy.GetMq(),
		Ex:      strategy.GetEx(),
		Created: time.Now(),
	}
	_, err := stratdb.InsertOne(mongoCtx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	return &proto.StoreStratRes{Strategy: strategy}, nil
}

type StratItem struct {
	Name    string    `bson:"_name,omitempty"`
	Mq      string    `bson:"mq"`
	Ex      string    `bson:"ex"`
	Created time.Time `bson:"created"`
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
	srv := &StratServiceServer{}
	proto.RegisterStratServiceServer(s, srv)
	reflection.Register(s)

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Initialize MongoDb client
	fmt.Println("Connecting to MongoDB...")

	// non-nil empty context
	mongoCtx = context.Background()

	// Connect takes in a context and options, the connection URI is the only option we pass for now
	// mongodb+srv://stockbrood:admin@stockbrood.sifn3lq.mongodb.net/test
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb+srv://stockbrood:admin@stockbrood.sifn3lq.mongodb.net/test"))
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

	c := make(chan os.Signal)
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
