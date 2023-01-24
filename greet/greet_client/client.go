package main

import (
	"GRPC-GO_COURSE/greet/greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"time"
)

func main() {
	fmt.Println(" HELLO Im Client")

	dialer, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer dialer.Close()

	client := greetpb.NewGreetServiceClient(dialer)

	//doUnary(client)
	//doServerStreaming(client)
	//doCLientStreaming(client)
	doBiDiStreaming(client)
}

func doUnary(client greetpb.GreetServiceClient) {
	fmt.Println("DOING A UNARY")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "ali",
			LastName:  "Forootan"},
	}
	res, err := client.Greet(context.Background(), req)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Resp: ", res.Result)
}

func doServerStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("DOING A SERVER STREAMING")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "ali",
			LastName:  "Forootan"},
	}
	res, err := client.GreetManyTimes(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 10; i++ {
		msg, err := res.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(msg.GetResult())
	}
}

func doCLientStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("DOING A Client STREAMING")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ali",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mamad",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "john",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "shima",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "lucy",
			},
		},
	}

	stream, err := client.LongGreet(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}

func doBiDiStreaming(client greetpb.GreetServiceClient) {
	fmt.Println("DOING A BiDi STREAMING")

	stream, err := client.GreetEveryOne(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	waitch := make(chan struct{})

	requests := []*greetpb.GreetEveryOneRequest{
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "ali",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mamad",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "john",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "shima",
			},
		},
		&greetpb.GreetEveryOneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "lucy",
			},
		},
	}

	go func() {
		for _, req := range requests {
			fmt.Println("Sending ", req.Greeting.FirstName)
			stream.Send(req)
			time.Sleep(1 * time.Second)
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
				fmt.Println(err)
				break
			}

			fmt.Println(" Received ", res.GetResult())
		}
		close(waitch)
	}()

	<-waitch
}
