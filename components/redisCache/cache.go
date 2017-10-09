package redisCache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"skylib/app"
	"time"
)

var Client *redis.Client

var redisOptions redis.Options

func connectToRedis() {
	Client = redis.NewClient(&redisOptions)
}

func Init() {
	if Client != nil {
		return
	}
	err := app.InitAppConfig()
	if err != nil {
		fmt.Println("config json not found")
		log.Println("config json not found")
		panic(err)
	}
	configRedis, err := app.GetConfigSelection("redis")
	if err != nil {
		fmt.Println("Redis not found in config")
		log.Println("Redis not found in config")
		return
	}

	var ok bool
	redisOptions.Addr, ok = configRedis["addr"].(string)
	if ok == false {
		fmt.Println("Redis addr not found in config")
		log.Println("Redis addr not found in config")
		panic(err)
	}

	connectToRedis()

	_, errPing := Client.Ping().Result()
	if errPing != nil {
		fmt.Println("Redis down", errPing.Error())
		log.Println("Redis down", errPing.Error())
		Client = nil
	}
	// Output: PONG <nil>
}

func Set(namespace string, key string, value interface{}) {
	if Client == nil {
		return
	}

	jsonM, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}

	err = Client.Set(namespace+":"+key, jsonM, 0).Err()
	if err != nil {
		log.Println(err)
	}
}

func Get(namespace string, key string) (string, error) {

	if Client == nil {
		return "", nil
	}

	val, err := Client.Get(namespace + ":" + key).Result()
	if err != nil && err != redis.Nil { // key not exists
		log.Println(err.Error())
		return "", err
	}
	return val, err

	//var bodyJson map[string]interface{}
	//err = json.Unmarshal([]byte(val), &bodyJson)
	//if err != nil {
	//	log.Println(err)
	//	return "", err
	//}
	//return bodyJson, nil
	// :-)))
}

func DeleteCache(key string) error {
	if Client == nil {
		return errors.New("cache not connet")
	}
	err := Client.Del(key).Err()
	if err != nil {
		log.Println(err)
	}
	return err
}

func GetCache(key string) (interface{}, error) {
	if Client == nil {
		return nil, errors.New("Client cache not connect")
	}

	result, err := Client.Get(key).Result()

	if err == redis.Nil {
		return nil, err
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var encJson map[string]interface{}
	err = json.Unmarshal([]byte(result), &encJson)

	return encJson, err
}

func SetCache(key string, value interface{}, expirations time.Duration) error {
	if Client == nil {
		return errors.New("cache not connet")
	}

	jsonM, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return err
	}

	err = Client.Set(key, jsonM, expirations).Err()
	if err != nil {
		log.Println("SetCache", err.Error())
	}
	return err
}
