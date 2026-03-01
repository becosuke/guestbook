package post

type Serial string

func NewSerial(serial string) *Serial {
	s := Serial(serial)
	return &s
}

func (s *Serial) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}
