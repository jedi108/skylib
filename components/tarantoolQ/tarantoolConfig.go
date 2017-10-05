package tarantoolQ

import (
	"skylib/app"
	"fmt"
	"os"
	"github.com/tarantool/go-tarantool"
)

type queueConfig struct {
	IsTarantoolServer bool
}

var (
	DebugText   []string
	Errors      []string
	Messages    []string
	QueueServer queueConfig
)

func GetLastErrorMessage() string {
	return Messages[len(Messages)-1]
}

func InitFromConfigTarantoolQueue() {
	Messages = append(Messages, "connect to tarantool")
	err := app.InitAppConfig()
	if err != nil {
		fmt.Println("config json not found")
		os.Exit(404)
	}
	configTarantool, err := app.GetConfigSelection("tarantool")
	if err != nil {
		Messages = append(Messages, "Not find: tarantool")
		QueueServer.IsTarantoolServer = false
		return
	}
	configIsEnableQueue, ok := configTarantool["enableQueue"]
	if ok == false {
		Messages = append(Messages, "Not find: enableQueue")
		QueueServer.IsTarantoolServer = false
		return
	}
	if configIsEnableQueue == false {
		Messages = append(Messages, "enableQueue: false")
		QueueServer.IsTarantoolServer = false
		return
	}

	configServersQueue, ok := configTarantool["serversQueue"]
	if ok == false {
		Messages = append(Messages, "Not find: serversQueue")
		QueueServer.IsTarantoolServer = false
		return
	}

	for _, valuesInServer := range configServersQueue.([]interface{}) {

		configSettingOfTarantool := valuesInServer.(map[string]interface{})

		/**
			Init Tarantool server connection
		 */
		cnnLocalHost := Conn{
			ConnectionUrl: configSettingOfTarantool["host"].(string),
			Opts: tarantool.Opts{
				User: configSettingOfTarantool["user"].(string),
				Pass: configSettingOfTarantool["pass"].(string),
			},
		}

		for nameValueServer, queuesServer := range valuesInServer.(map[string]interface{}) {

			/** host, user, password  */
			switch queuesServer.(type) {
			case string:
				nameValueServer = nameValueServer //not use

			/** queues  */
			case []interface{}:
				for _, configListFieldsInSettingServer := range queuesServer.([]interface{}) {

					nameOfQueue := configListFieldsInSettingServer.(string)

					/**
						Init queue to pool queues
					 */
					err := QTarantool.PushToQueues(&QueueT{
						ConnectionTarantool: &cnnLocalHost,
						QueueName:           nameOfQueue,
					})

					if err != nil {
						Messages = append(Messages, err.Error())
						QueueServer.IsTarantoolServer = false
						return
					}

					DebugText = append(DebugText, "--"+fmt.Sprintln(configListFieldsInSettingServer))
				}
			}

		}
	}
	QueueServer.IsTarantoolServer = true
}
