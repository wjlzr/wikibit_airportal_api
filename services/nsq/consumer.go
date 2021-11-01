package nsq

import (
	"strings"
	"time"
	"wiki_bit/boot/log"

	"go.uber.org/zap"

	nsq "github.com/nsqio/go-nsq"
)

//初始化消费者
func (m Nsq) InitConsumer(topic, channel string, handler nsq.Handler) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 1 * time.Second      //设置重连时间
	c, err := nsq.NewConsumer(topic, channel, cfg) // 新建一个消费者
	if err != nil {
		log.Logger().Error("Nsq create consumer err: ", zap.Error(err))
	}
	c.AddHandler(handler) // 添加消费者接口

	//建立NSQLookupd连接
	addr := strings.Builder{}
	addr.WriteString(m.Host)
	addr.WriteString(":")
	addr.WriteString(m.Port)
	if err := c.ConnectToNSQD(addr.String()); err != nil {
		log.Logger().Error("Nsq ConnectToNSQD err: ", zap.Error(err))
		return
	}
	return
	//<-c.StopChan
}
