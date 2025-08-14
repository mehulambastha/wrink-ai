package services

import (
	"context"
	"log"
	"time"

	pb "github.com/mehulambastha/wrink-ai/backend/gen/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcSearchService struct {
	client pb.SearchServiceClient
}

func NewGRPCSearchService(serverAddr string) (*grpcSearchService, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	log.Printf("Succesfully connected to gRPC search service at %s", serverAddr)

	client := pb.NewSearchServiceClient(conn)

	return &grpcSearchService{client: client}, nil
}

func (s *grpcSearchService) Search(ctx context.Context, keywords string, topic string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)

	defer cancel()

	req := &pb.SearchRequest{
		Topic:    topic,
		Keywords: []string{keywords},
	}

	res, err := s.client.Search(ctx, req)

	if err != nil {
		log.Printf("gRPC search call failed: %v", err)
		return "", err
	}

	return res.GetSearchResultsText(), nil
}
