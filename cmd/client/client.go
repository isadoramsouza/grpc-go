package main

import (
	"context"
	"github.com/isadoramsouza/grpc-go/pb"
	"google.golang.org/grpc"
	"log"
)

func main(){
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	AddUser(client)
}

func AddUser(client pb.UserServiceClient)  {
	req := &pb.User{
		Id: "1",
		Name: "Isa",
		Email: "isa@teste.com",
	}
	res,err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}
	log.Println(res)
}