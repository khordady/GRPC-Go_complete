package main

import (
	"GRPC-GO_COURSE/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
	"time"
)

type server struct {
	greetpb.GreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("GREET INVOKED \n", req)
	first := req.GetGreeting().GetFirstName()
	result := "Hello " + first
	res := &greetpb.GreetResponse{Result: result}

	return res, nil
}

func (*server) GreettManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GREET INVOKED \n", req)
	first := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + first + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{Result: result}

		stream.Send(res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Long GREET INVOKED \n ", stream)
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			fmt.Println(err)
		}

		firstname := req.GetGreeting().GetFirstName()
		result += "Hello " + firstname + "! \n "
	}
}

func (*server) GreetEveryOne(stream greetpb.GreetService_GreetEveryOneServer) error {
	fmt.Println("EveryOne GREET INVOKED \n ")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println(err)
		}

		firstname := req.GetGreeting().GetFirstName()
		result := "Hello " + firstname + "! \n "

		err = stream.Send(&greetpb.GreetEveryOneResponse{
			Result: result,
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		fmt.Println(err)
		return
	}

	ser := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(ser, &server{})

	if err := ser.Serve(lis); err != nil {
		fmt.Println(err)
		return
	}
}
