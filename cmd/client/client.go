package main

import (
	"context"
	"fmt"
	"github.com/isadoramsouza/grpc-go/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "1",
		Name:  "Isa",
		Email: "isa@teste.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}
	log.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "1",
		Name:  "Isa",
		Email: "isa@teste.com",
	}
	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive message %v", err)
		}
		println("Status ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Isa",
			Email: "teste@teste.com",
		},
		{
			Id:    "2",
			Name:  "Fulano",
			Email: "teste@teste.com",
		},
		{
			Id:    "3",
			Name:  "Ciclano",
			Email: "teste@teste.com",
		},
		{
			Id:    "4",
			Name:  "Beltrano",
			Email: "teste@teste.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not create request %v", err)
	}
	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Could not receive response %v", err)
	}
	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Could not create request %v", err)
	}

	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Isa",
			Email: "teste@teste.com",
		},
		{
			Id:    "2",
			Name:  "Fulano",
			Email: "teste@teste.com",
		},
		{
			Id:    "3",
			Name:  "Ciclano",
			Email: "teste@teste.com",
		},
		{
			Id:    "4",
			Name:  "Beltrano",
			Email: "teste@teste.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending User ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 3)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not receive stream %v", err)
				break
			}
			fmt.Printf("Receiving user %v with status %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
