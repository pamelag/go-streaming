package handler

import (
	"errors"
	"fmt"
	"os"

	uuid "github.com/satori/go.uuid"
)

const (
	logPath = "/home/pamela/go/src/github.com/pamelag/go-streaming/logs"

	READY  = 0
	ACTIVE = 1
	FULL   = -1

	LOCKED = 0
	OPEN   = 1

	BUFFER1 = 0
	BUFFER2 = 1
)

/** Segment **/
var segment *Segment

/** Path **/
var Path = "/home/pamela/go/src/github.com/pamelag/go-streaming/logs/"

/** createTopic is used to create the topic across
the cluster. **/

func CreateTopic(name string) (*Topic, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	t := &Topic{
		ID:   uid.String(),
		Name: name,
	}
	return t, nil
}

/**

BuildPartition

**/
func BuildPartition(t *Topic, p string, user string) (*Partition, error) {

	partition, exists := t.ActivePartions[p]

	if exists || partition.Status == ACTIVE {
		return nil, errors.New("Partition exists or in an active mode")
	}

	pID, _ := uuid.NewV1()

	partition = &Partition{
		ID:              pID.String(),
		Name:            p,
		PartitionHandle: partition,
		Topic:           t.Name,
		MaxBuffers:      2,
	}

	for i := 0; i < partition.MaxBuffers; i++ {
		stream := &Stream{
			ID:       i,
			Status:   READY,
			Lock:     OPEN,
			Size:     0,
			Messages: make([]Message, 500),
			MaxSize:  500,
		}
		partition.Streams = append(partition.Streams, stream)
	}
	err1 := CreatePartionOnDisk(partition)

	if err1 != nil {
		return nil, err1
	}
	return partition, nil
}

/**
GetActiveStream
*/
func GetActiveStream(p *Partition) (*Stream, error) {
	stream1, stream2 := p.Streams[0], p.Streams[1]

	switch {
	case stream1.Status == LOCKED && stream2.Status == LOCKED:
		return nil, errors.New("buffers not available for the new message")
	case stream1.Status == READY && stream2.Status == READY:
		return stream1, nil
	case stream1.Status == ACTIVE:
		return stream1, nil
	case stream2.Status == ACTIVE:
		return stream2, nil
	default:
		return nil, nil
	}
}

func Write(stream *Stream, m Message, activeBuffer chan int, buffStat chan int) {
	fmt.Printf("Switching to %s ", stream.Name)
	messages := append(stream.Messages, m)
	stream.Size = len(messages)
	// from here
	if stream.Size < 500 {
		activeBuffer <- stream.ID
		stream.Status = ACTIVE
	} else {
		stream.Status = FULL
		buffStat <- FULL
		if stream.ID == 0 {
			activeBuffer <- -1
		} else {
			activeBuffer <- 0
		}
	}
}

/**  CreatePartionOnDisk  **/

func CreatePartionOnDisk(prt *Partition) error {
	prt.Path = Path + prt.Name
	err := os.MkdirAll(prt.Name, 0755)
	if err != nil {
		return err
	}
	return nil
}

/**  WriteToLog  **/

func WriteToLog(stream *Stream, status int, timestamp string) {
	if status == -1 {
		stream.Status = LOCKED
		path := Path + stream.Name + "_" + timestamp + "_" + ".log"
		segment := NewSegment(stream.Name, path)
		segment.CommitToFile(stream.Messages)
	}
}

// c := make(chan int)
// done := make(chan bool)

// go func() {
// 	for i := 0; i < 10; i++ {
// 		c <- i
// 	}
// 	done <- true
// }()

// go func() {
// 	<-done
// 	close(c)
// }()

// for n := range c {
// 	fmt.Println(n)
// }
