package mocks

type PingerStub struct{}

func NewPingerStub() *PingerStub {
	return &PingerStub{}
}

func (p *PingerStub) Ping() error {
	return nil
}
