package grpcapi

import (
	"context"
	"fmt"
	pb "github.com/sunwild/domain-checker_checker/proto/checker/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func CheckDomainsWithGRPC(domains []string, isManual bool) (*pb.DomainResponse, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDomainCheckerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	req := &pb.DomainRequest{
		Domains:  domains,
		IsManual: isManual,
	}

	response, err := client.CheckDomains(ctx, req)
	if err != nil {
		log.Printf("gRPC request failed: %v", err)
		return nil, fmt.Errorf("gRPC request failed: %w", err)
	}

	return response, nil
}
