package im

type Message struct {
	Head struct {
		OpId uint64
		Type string
		Path string
	}
	Body interface{}
}
