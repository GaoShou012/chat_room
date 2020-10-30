package room

import (
	"encoding/json"
	"runtime"
	proto_room "wchatv1/proto/room"
)

type decodeJob struct {
	input  []byte
	output *proto_room.Message
	done   chan bool
	err    error
}

type encodeJob struct {
	input  *proto_room.Message
	output []byte
	done   chan bool
	err    error
}

type codec struct {
	decodeJobs chan *decodeJob
	encodeJobs chan *encodeJob
}

func (c *codec) init(jobsCacheSize int) {
	c.decodeJobs = make(chan *decodeJob, jobsCacheSize)
	c.encodeJobs = make(chan *encodeJob, jobsCacheSize)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(c *codec) {
			for {
				select {
				case job := <-c.decodeJobs:
					job.err = json.Unmarshal(job.input, job.output)
					job.done <- true
					break
				case job := <-c.encodeJobs:
					job.output, job.err = json.Marshal(job.input)
					job.done <- true
					break
				}
			}
		}(c)
	}
}

func (c *codec) Decode(input []byte) (*proto_room.Message, error) {
	job := decodeJob{
		input:  input,
		output: &proto_room.Message{},
		done:   make(chan bool, 1),
	}
	c.decodeJobs <- &job
	<-job.done
	return job.output, job.err
}

func (c *codec) Encode(message *proto_room.Message) ([]byte, error) {
	job := encodeJob{
		input:  message,
		output: nil,
		done:   make(chan bool, 1),
		err:    nil,
	}
	c.encodeJobs <- &job
	<-job.done
	return job.output, job.err
}
