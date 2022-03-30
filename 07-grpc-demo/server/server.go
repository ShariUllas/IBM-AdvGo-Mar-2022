package main

import (
	"context"
	"fmt"
	"grpc-demo/proto"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type appService struct {
	proto.UnimplementedAppServiceServer
}

func (s *appService) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
	x := req.GetX()
	y := req.GetY()
	result := x + y
	fmt.Printf("Add Operation: X = %d, Y = %d\n", x, y)
	res := &proto.AddResponse{
		Result: result,
	}
	return res, nil
}

func (s *appService) GeneratePrimes(req *proto.PrimeRequest, serverStream proto.AppService_GeneratePrimesServer) error {
	start := req.GetStart()
	end := req.GetEnd()
	for no := start; no <= end; no++ {
		if isPrime(no) {
			time.Sleep(500 * time.Millisecond)
			res := &proto.PrimeResponse{
				PrimeNo: no,
			}
			fmt.Printf("Sending Prime No : %d\n", no)
			serverStream.Send(res)
		}
	}
	fmt.Println("All Prime numbers are sent!")
	return nil
}

func isPrime(no int32) bool {
	for i := int32(2); i <= (no / 2); i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}

func (s *appService) CalculateAverage(serverStream proto.AppService_CalculateAverageServer) error {
	var count int32
	var sum int32

	for {
		req, err := serverStream.Recv()
		if err == io.EOF {
			avg := sum / count
			res := &proto.AverageResponse{
				Count:   count,
				Average: avg,
			}
			serverStream.SendAndClose(res)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Average Request: no = ", req.GetNo())
		count++
		sum += req.GetNo()
	}
	return nil
}

func main() {
	s := &appService{}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAppServiceServer(grpcServer, s)
	grpcServer.Serve(listener)
}
