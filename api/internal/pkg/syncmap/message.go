package syncmap

type Message struct {
	serial int64
	body   string
}

func NewMessage(serial int64, body string) *Message {
	return &Message{
		serial: serial,
		body:   body,
	}
}

func (m Message) Serial() int64 {
	return m.serial
}

func (m Message) Body() string {
	return m.body
}
