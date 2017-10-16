package tarantoolQ

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tarantool/go-tarantool"
	"github.com/tarantool/go-tarantool/queue"
	"log"
	"os"
	"skylib/app"
	"time"
)

var SkyQ SkyQueue

type SkyQueue struct {
	Queue         queue.Queue
	host          string
	user          string
	pass          string
	Err           error
	queueName     string
	IsEnableQueue bool
	StatTasks     map[string]int
	StatCalls     map[string]int
}

func NewConnect() SkyQueue {
	if SkyQ.IsEnableQueue == false {
		SkyQ.Err = initTarantoolQueue()
		if SkyQ.Err == nil {
			SkyQ.ConnectTarantool()
		}
	}
	return SkyQ
}

func initTarantoolQueue() error {
	err := app.InitAppConfig()
	if err != nil {
		log.Println(err)
		fmt.Println("config json not found")
		os.Exit(404)
	}
	configTarantool, err := app.GetConfigSelection("tarantool2")
	if err != nil {
		log.Println(err)
		return errors.New("Not find: tarantool")
	}
	configIsEnableQueue, ok := configTarantool["enableQueue"]
	if ok == false {
		return errors.New("Not find: enableQueue")
	}
	if configIsEnableQueue == false {
		return errors.New("enableQueue: FALSE")
	}
	SkyQ.host, ok = configTarantool["host"].(string)
	if ok == false {
		log.Println(err)
		return errors.New("Not find: host")
	}
	SkyQ.user, ok = configTarantool["user"].(string)
	if ok == false {
		log.Println(err)
		return errors.New("Not find: user")
	}
	SkyQ.pass, ok = configTarantool["pass"].(string)
	if ok == false {
		log.Println(err)
		return errors.New("Not find: pass")
	}
	SkyQ.queueName, ok = configTarantool["queue"].(string)
	if ok == false {
		log.Println(err)
		return errors.New("Not find: queue")
	}
	//SkyQ.IsEnableQueue = true

	return nil
}

func (skyQueue *SkyQueue) ParseStatistic() error {
	stat, err := skyQueue.Queue.Statistic()
	if err != nil {
		return err
	}
	v := stat.(map[interface{}]interface{})
	for i, s := range v {
		if i == "tasks" {
			err = mapstructure.Decode(s, &skyQueue.StatTasks)
			if err != nil {
				return err
			}
		}
		if i == "calls" {
			err := mapstructure.Decode(s, &skyQueue.StatCalls)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (skyQueue *SkyQueue) EchoStatistics() error {
	err := skyQueue.ParseStatistic()
	if err != nil {
		return err
	}
	fmt.Print("  _ buried:", skyQueue.StatTasks["buried"])
	fmt.Print("  _ delayed:", skyQueue.StatTasks["delayed"])
	fmt.Print("  _ done:", skyQueue.StatTasks["done"])
	fmt.Print("  _ ready:", skyQueue.StatTasks["ready"])
	fmt.Print("  _ taken:", skyQueue.StatTasks["taken"])
	fmt.Print("  _ total:", skyQueue.StatTasks["total"])
	fmt.Println()
	return nil
}

func t() time.Duration {
	return time.Second * 2
}

func (skyQueue *SkyQueue) ConnectTarantool() SkyQueue {
	//cnn, err := tarantool.Connect(skyQueue.host, tarantool.Opts{User: skyQueue.user, Pass: skyQueue.pass, Timeout:time.Second * 2, MaxReconnects:3})
	cnn, err := tarantool.Connect(skyQueue.host, tarantool.Opts{User: skyQueue.user, Pass: skyQueue.pass})
	if err != nil {
		skyQueue.Err = err
		log.Println(err)
		return *skyQueue
	}
	skyQueue.Queue = queue.New(cnn, skyQueue.queueName)
	_, errEx := skyQueue.Queue.Exists()
	if err != nil {
		skyQueue.Err = errEx
		log.Println(err)
		return *skyQueue
	}
	skyQueue.IsEnableQueue = true
	return *skyQueue
}
