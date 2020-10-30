package room

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"runtime"
	proto_room "wchatv1/proto/room"
)

type job struct {
	done chan bool
	err  error
}

type encodeMessageJob struct {
	job
	input  *proto_room.Message
	output []byte
}
type decodeMessageJob struct {
	job
	input  []byte
	output *proto_room.Message
}

type encodePassTokenJob struct {
	job
	key    []byte
	input  *proto_room.PassToken
	output string
}
type decodePassTokenJob struct {
	job
	key    []byte
	input  string
	output *proto_room.PassToken
}

type codec struct {
	encodeMessageJob chan *encodeMessageJob
	decodeMessageJob chan *decodeMessageJob

	encodePassTokenJob chan *encodePassTokenJob
	decodePassTokenJob chan *decodePassTokenJob
}

func (c *codec) Init() {
	fmt.Println("Init Codec")

	c.encodeMessageJob = make(chan *encodeMessageJob, runtime.NumCPU()*1000)
	c.decodeMessageJob = make(chan *decodeMessageJob, runtime.NumCPU()*1000)

	c.encodePassTokenJob = make(chan *encodePassTokenJob, runtime.NumCPU()*1000)
	c.decodePassTokenJob = make(chan *decodePassTokenJob, runtime.NumCPU()*1000)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(c *codec) {
			for {
				select {
				case job := <-c.encodePassTokenJob:
					j, err := json.Marshal(job.input)
					if err != nil {
						job.err = err
						job.done <- true
						break
					}
					m := jwt.MapClaims{}
					if err := json.Unmarshal(j, &m); err != nil {
						job.err = err
						job.done <- true
						break
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, m)
					str, err := token.SignedString(job.key)
					if err != nil {
						job.err = err
						job.done <- true
						break
					}
					job.output, job.err = str, nil
					job.done <- true
				case job := <-c.decodePassTokenJob:
					token, err := jwt.Parse(job.input, func(token *jwt.Token) (i interface{}, err error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("Unexpected signing method: %v\n", token.Header["alg"])
						}
						return job.key, nil
					})
					if err != nil {
						job.err = err
						job.done <- true
						break
					}

					passToken := &proto_room.PassToken{}
					if err := mapstructure.WeakDecode(token.Claims.(jwt.MapClaims), passToken); err != nil {
						job.err = err
						job.done <- true
						break
					}

					job.output, job.err = passToken, nil
					job.done <- true
				case job := <-c.encodeMessageJob:
					j, err := json.Marshal(job.input)
					if err != nil {
						job.err = err
						job.done <- true
						break
					}
					job.output, job.err = j, nil
					job.done <- true
				case job := <-c.decodeMessageJob:
					msg := &proto_room.Message{}
					if err := json.Unmarshal(job.input, msg); err != nil {
						job.err = err
						job.done <- true
						break
					}
					job.output, job.err = msg, nil
					job.done <- true
				}
			}
		}(c)
	}
}

func (c *codec) EncodeMessage(input *proto_room.Message) ([]byte, error) {
	job := &encodeMessageJob{
		job:    job{done: make(chan bool, 1), err: nil},
		input:  input,
		output: nil,
	}
	c.encodeMessageJob <- job
	<-job.done
	return job.output, job.err
}
func (c *codec) DecodeMessage(input []byte) (*proto_room.Message, error) {
	job := &decodeMessageJob{
		job:    job{done: make(chan bool, 1), err: nil},
		input:  input,
		output: nil,
	}
	c.decodeMessageJob <- job
	<-job.done
	return job.output, job.err
}

func (c *codec) EncodePassToken(key []byte, input *proto_room.PassToken) (string, error) {
	job := &encodePassTokenJob{
		job:    job{done: make(chan bool, 1), err: nil},
		key:    key,
		input:  input,
		output: "",
	}
	c.encodePassTokenJob <- job
	<-job.done
	return job.output, job.err
}
func (c *codec) DecodePassToken(key []byte, input string) (*proto_room.PassToken, error) {
	job := &decodePassTokenJob{
		job:    job{done: make(chan bool, 1), err: nil},
		key:    key,
		input:  input,
		output: nil,
	}
	c.decodePassTokenJob <- job
	<-job.done
	return job.output, job.err
}
