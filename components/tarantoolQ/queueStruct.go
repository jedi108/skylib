package tarantoolQ

import (
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
	"log"
)

var QTarantool QueuesTarantool

//func init() {
//	QTarantool = *new(QueuesTarantool)
//}

type Conn struct {
	ConnectionUrl string
	Opts          tarantool.Opts
	connect       *tarantool.Connection
}

func (conn *Conn) ConnectTaran() error {
	connq, err := tarantool.Connect(conn.ConnectionUrl, conn.Opts)
	conn.connect = connq
	return err
}

type QueueT struct {
	ConnectionTarantool *Conn
	QueueName           string
	queue               *queue.Queue
	QueueCfg            queue.Cfg
}

func (queueT *QueueT) connectQueue() {
	queueVal := queue.New(queueT.ConnectionTarantool.connect, queueT.QueueName)
	queueT.queue = &queueVal
}

//for creating new queue, now not using
func (queueT *QueueT) ConnectCreateQueue() (error) {
	queueVal := queue.New(queueT.ConnectionTarantool.connect, queueT.QueueName)
	queueT.queue = &queueVal
	return queueVal.Create(queueT.QueueCfg)
}

type QueuesTarantool struct {
	Queues map[string]*QueueT
}

func (queues *QueuesTarantool) SetQueue(nameQ string, queueT *QueueT) {
	if len(queues.Queues) == 0 {
		queues.Queues = make(map[string]*QueueT)
	}
	queues.Queues[nameQ] = queueT
}

func (queues *QueuesTarantool) GetQueue(nameQueue string) *queue.Queue {
	return queues.Queues[nameQueue].queue
}

func (qt *QueuesTarantool) PushToQueues(queueT *QueueT) {

	_, ok := qt.Queues[queueT.QueueName]
	if true == ok {
		log.Println("Queue `%s` already added", queueT.QueueName)
	}

	err := queueT.ConnectionTarantool.ConnectTaran()
	if err != nil {
		log.Println("Error connect tarantool:")
		log.Println("Queue opts", queueT.ConnectionTarantool.Opts)
		log.Println("ConnectionUrl", queueT.ConnectionTarantool.ConnectionUrl)
		log.Println(err.Error())
		panic(err)
	}

	queueT.connectQueue()

	qt.SetQueue(queueT.QueueName, queueT)

	//qt.Queues[queueT.QueueName] = queueT
}
