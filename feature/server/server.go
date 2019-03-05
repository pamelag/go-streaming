package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/pamelag/go-streaming/feature/featurepb"
	"github.com/pamelag/go-streaming/handler"
	"google.golang.org/grpc"
)

var path string
var topics map[string]*handler.Topic
var activePartitions map[string]*handler.Partition

//var activePrt chan *Partition
var activeBuff chan int
var activeBuffState chan int

var Path string = "/home/pamela/go/src/github.com/pamelag/go-streaming/logs/"

const (
	logPath = "/home/pamela/go/src/github.com/pamelag/go-streaming/logs"
	READY   = 0
	ACTIVE  = 1
	FULL    = -1

	LOCKED = 0
	OPEN   = 1

	BUFFER1 = 0
	BUFFER2 = 1

	UNREAD = 0
	READ   = 1
)

func init() {
	topicNames := []string{"Document", "Wireframe", "Image"}
	for _, name := range topicNames {
		tp, err := handler.CreateTopic(name)
		if err != nil {
			fmt.Printf("Error %v ", err)
		}
		topics[name] = tp
	}

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

type server struct {
}

func (*server) CreateDocument(ctx context.Context, req *featurepb.DocumentRequest) (*featurepb.DocumentResponse, error) {
	fmt.Println("CreateDocument function was invoked with %v", req)
	p := req.GetDocument().GetName()
	t := req.GetDocument().GetTopic()
	user := req.GetDocument().GetUser()
	tm := req.GetDocument().GetTimestamp()

	topic := topics[t]

	partition, err := handler.BuildPartition(topic, p, user)
	if err != nil {
		return nil, err
	}

	result := " Document Created " + partition.ID + " : " + partition.Name
	res := &featurepb.DocumentResponse{
		Result: result,
	}
	//add to active partitions
	activePartitions[p] = partition
	return res, nil
}

func (*server) FeatureUpdates(ctx context.Context, req *featurepb.FeatureRequest) (*featurepb.FeatureResponse, error) {
	fmt.Println("FeatureUpdates function was invoked with %v", req)
	featureName := req.GetFeature().GetName()
	result := "Updating " + featureName + " " + req.GetFeature().GetDelta()
	res := &featurepb.FeatureResponse{
		Result: result,
	}
	return res, nil
}

func (*server) LongFeatureUpdates(stream featurepb.FeatureService_LongFeatureUpdatesServer) error {
	fmt.Printf("LongFeatureUpdates function was invoked with a streamimg request")
	result := ""
	activeBuff = make(chan int)
	activeBuffState = make(chan int)
	prt := &handler.Partition{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			bufferStream, err := handler.GetActiveStream(prt)
			if err != nil {
				return err
			}
			bufferStream.Status = FULL

			// we have finished the client stream
			return stream.SendAndClose(&featurepb.LongFeatureResponse{
				Status: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		name := req.GetFeature().GetName()
		topic := req.GetFeature().GetTopic()
		user := req.GetFeature().GetUser()
		delta := req.GetFeature().GetDelta()
		operation := req.GetFeature().GetOperation()
		timeStamp := req.GetFeature().GetTimestamp()

		message := handler.Message{
			Topic:     topic,
			Partition: name,
			User:      user,
			Delta:     delta,
			Operation: operation,
			Timestamp: timeStamp,
		}

		prt = activePartitions[name]
		bufferStream, err := handler.GetActiveStream(prt)

		if err != nil {
			return err
		}
		go handler.Write(bufferStream, message, activeBuff, activeBuffState)

		actBuff := <-activeBuff
		bufstat := <-activeBuffState

		fmt.Println("Active buffer id and state ", actBuff, bufstat)

		go handler.WriteToLog(bufferStream, bufstat, timeStamp)
		result += " " + delta + " "
	}
}

func main() {
	fmt.Println("Hello world")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("failed to listen %v ", err)
	}

	s := grpc.NewServer()
	featurepb.RegisterFeatureServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
