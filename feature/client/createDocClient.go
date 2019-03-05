package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pamelag/go-streaming/feature/featurepb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer cc.Close()
	c := featurepb.NewFeatureServiceClient(cc)

	createDocClient(c)
}

func createDocClient(c featurepb.FeatureServiceClient) {
	fmt.Println("Create Client Doc")
	req := &featurepb.DocumentRequest{
		Document: &featurepb.Document{
			Name:      "Collecting Metrices",
			Topic:     "Feature",
			Timestamp: "2019-10-20 00:00:00.000 ",
			User:      "Tics",
		},
	}
	res, err := c.CreateDocument(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Create Document: %v", err)
	}
	fmt.Printf("Response from Feature Request %v ", res.Result)
}
