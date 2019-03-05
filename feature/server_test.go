package server_test

import (
	"context"
	"fmt"
	"log"

	"github.com/pamelag/go-streaming/feature/featurepb"
	"github.com/pamelag/jocko/jocko"
	"google.golang.org/grpc"
	//"github.com/pamelag/go-streaming/feature/server"
)

const (
	topic = "test_topic"
)

func init() {
	log.SetLevel("debug")
}

func createDocument(t ti.T, s1 *jocko.Server, other ...*jocko.Server) error {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer cc.Close()
	c := featurepb.NewFeatureServiceClient(cc)

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
