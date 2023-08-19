package syncmap

type Message struct {
	key   string
	value string
}

func NewMessage(key string, value string) *Message {
	return &Message{
		key:   key,
		value: value,
	}
}

func (m Message) Key() string {
	return m.key
}

func (m Message) Value() string {
	return m.value
}
