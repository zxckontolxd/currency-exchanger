package main

import (
	"log"
	"net"
	pb "github.com/zxckontolxd/proto-exchange/exchange"
	"context"
    "google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExchangeServiceServer
}

func (s *server) GetExchangeRates (ctx context.Context, _ *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	// REWRITE! Need db. 
	// TESTCODE
	rates := map[string]float32{
        "USD": 1.0,
        "EUR": 0.9,
        "RUB": 90.0,
    }
	return &pb.ExchangeRatesResponse{Rates: rates}, nil
}

func (s *server) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	//TESTCODE!!!
	var rate float32
	rate = 0.9
	return &pb.ExchangeRateResponse{
        FromCurrency: req.FromCurrency,
        ToCurrency:   req.ToCurrency,
        Rate:         rate,
    }, nil
}

func main () {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Filed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, &server{})
	log.Printf("Server listen at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Field to serve: %v", err)
	}
}