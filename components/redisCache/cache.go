package redisCache

import (
	"github.com/go-redis/redis"
	"log"
	"encoding/json"
	"encoding/hex"
	"errors"
	"fmt"
)

var Client *redis.Client

type RedisSource struct {
	Addr     string
	Password string
}

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := Client.Ping().Result()
	if err != nil {
		log.Println("Redis down")
	}
	// Output: PONG <nil>
}

func Set(namespace string, key string, value interface{}) {
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

func GetStruct(namespace string, key []byte) (interface{}, error) {
	if Client == nil {
		return nil, errors.New("Client cache not connect")
	}
	result, err := Client.Get(namespace + ":" + hex.EncodeToString(key)).Result()

	fmt.Println("_______________________")
	fmt.Println(result)

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var encJson map[string]interface{}
	err = json.Unmarshal([]byte(result), &encJson)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return encJson, nil
}

func SetStruct(namespace string, key []byte, value interface{}) {
	if Client == nil {
		return
	}
	nameSpaceKey := hex.EncodeToString(key)

	jsonM, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return
	}

	err = Client.Set(namespace+":"+nameSpaceKey, jsonM, 0).Err()
	if err != nil {
		log.Println(err)
	}
}
