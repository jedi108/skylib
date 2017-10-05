package tarantoolQ

import (
	"github.com/mitchellh/mapstructure"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
	"log"
	"os"
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
	StatTasks           map[string]int
	StatCalls           map[string]int
}

func NewWorker(queueString string) QueueT {
	if QueueServer.IsTarantoolServer == false {
		log.Println("config api tarantool if off")
		log.Println(Messages)
		os.Exit(1)
	}
	return *QTarantool.GetTQueue(queueString)
}

func (queues *QueuesTarantool) GetTQueue(nameQueue string) *QueueT {
	return queues.Queues[nameQueue]
}

func (queueT *QueueT) connectQueue() {
	queueVal := queue.New(queueT.ConnectionTarantool.connect, queueT.QueueName)
	queueT.queue = &queueVal
}

func (queueT *QueueT) GetQ() *queue.Queue {
	return queueT.queue
}

func (queueT *QueueT) ParseStatistic() error {
	qqq := *(queueT.queue)
	stat, err := qqq.Statistic()
	if err != nil {
		return err
	}
	v := stat.(map[interface{}]interface{})
	for i, s := range v {
		if i == "tasks" {
			err = mapstructure.Decode(s, &queueT.StatTasks)
			if err != nil {
				return err
			}
		}
		if i == "calls" {
			err := mapstructure.Decode(s, &queueT.StatCalls)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func (qt *QueuesTarantool) PushToQueues(queueT *QueueT) error {

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
		return err
	}

	queueT.connectQueue()

	qt.SetQueue(queueT.QueueName, queueT)

	return nil
}
