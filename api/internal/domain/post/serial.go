package post

type Serial int64

func NewSerial(serial int64) *Serial {
	s := Serial(serial)
	return &s
}

func (s *Serial) Int64() int64 {
	if s == nil {
		return 0
	}
	return int64(*s)
}
