package nsq

import (
	"wiki_bit/config"
)

type Nsq struct {
	Topic string
	Host  string
	Port  string
}

func NewNsq() Nsq {
	return Nsq{
		Host: config.Conf().Nsq.Host,
		Port: config.Conf().Nsq.Port,
	}
}
