package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/thg021/comunicacao-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connection to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	//Estamos adicionando nossos servicos
	// Linhas comentadas apenas para teste
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUsersStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Thiago",
		Email: "teste@g.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make  to gRPC Request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Thiago",
		Email: "teste@g.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make  to gRPC Request: %v", err)
	}

	for {
		stream, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could note recive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "t1",
			Name:  "Thiago S",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t2",
			Name:  "Heitor",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t3",
			Name:  "Mirella",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t4",
			Name:  "Debora",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t5",
			Name:  "Elias",
			Email: "t@t.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUsersStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUsersStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "t1",
			Name:  "Thiago S",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t2",
			Name:  "Heitor",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t3",
			Name:  "Mirella",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t4",
			Name:  "Debora",
			Email: "t@t.com",
		},
		&pb.User{
			Id:    "t5",
			Name:  "Elias",
			Email: "t@t.com",
		},
	}

	//Channel é um local onde voce manda uma comunicação entre rotrinas
	wait := make(chan int)

	//Go routines, ficara rodando em background
	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	//Em paralelo / concorrente

	go func() {
		for {

			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}

			fmt.Printf("Receiving user %v with status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}

		//fechando o channel
		close(wait)
	}()

	<-wait
}
