package handler

type Topic struct {
	ID             string
	Name           string
	ActivePartions map[string]*Partition
	numPartitions  []int
}

type Partition struct {
	ID              string
	Name            string
	PartitionHandle *Partition
	Topic           string
	Path            string
	Streams         []*Stream
	Status          int
	MaxBuffers      int
}

func (p *Partition) GetStreams() []*Stream {
	return p.Streams
}

type Stream struct {
	ID       int
	Name     string
	Status   int
	Lock     int
	Size     int
	Messages []Message
	MaxSize  int
}

type Message struct {
	Topic     string
	Partition string
	User      string
	Delta     string
	Operation string
	Timestamp string
}
