package room

import (
	"github.com/GaoShou012/frontier"
	"sync"
)

type senderParallelJob struct {
	conn frontier.Conn
	data []byte
}
type senderParallelConnId struct {
	maxProcess int
	jobsPool   sync.Pool
	jobs       []chan *senderParallelJob
}

func (s *senderParallelConnId) Init(maxProcess int) {
	s.maxProcess = maxProcess
	s.jobsPool.New = func() interface{} {
		return new(senderParallelJob)
	}
	s.jobs = make([]chan *senderParallelJob, maxProcess)
	for i := 0; i < maxProcess; i++ {
		s.jobs[i] = make(chan *senderParallelJob, 10000)
	}
	for i := 0; i < maxProcess; i++ {
		go func(i int) {
			for {
				job := <-s.jobs[i]
				conn, data := job.conn, job.data
				conn.Sender(data)
				//SenderMessageOnQueue.Dec()
				// release the job struct
				s.jobsPool.Put(job)
			}
		}(i)
	}
}

func (s *senderParallelConnId) Send(conn frontier.Conn, data []byte) {
	//SenderMessageOnQueue.Inc()

	// get the job struct from pool
	job := s.jobsPool.Get().(*senderParallelJob)
	job.conn = conn
	job.data = data

	index := conn.GetId() % s.maxProcess
	s.jobs[index] <- job
}
