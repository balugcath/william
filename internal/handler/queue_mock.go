package handler

type queueMock struct {
	res interface{}
}

func (s *queueMock) Add(r interface{}) {
	s.res = r
}
