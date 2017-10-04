package redisCache

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"log"
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

func DeleteCache(key string) error {
	if Client == nil {
		return errors.New("cache not connet")
	}
	return Client.Del(key).Err()
}

func GetCache(key string) (interface{}, error) {
	if Client == nil {
		return nil, errors.New("Client cache not connect")
	}
	log.Println("key GetCache", key)
	result, err := Client.Get(key).Result()

	if err == redis.Nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var encJson map[string]interface{}
	err = json.Unmarshal([]byte(result), &encJson)

	return encJson, err
}

func SetCache(key string, value interface{}) error {
	if Client == nil {
		return errors.New("cache not connet")
	}

	jsonM, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Client.Set(key, jsonM, 0).Err()
}

/**
	DEPRECATED
 */
func GetStruct(namespace string, key []byte) (interface{}, error) {
	if Client == nil {
		return nil, errors.New("Client cache not connect")
	}
	result, err := Client.Get(namespace + ":" + hex.EncodeToString(key)).Result()

	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var encJson map[string]interface{}
	err = json.Unmarshal([]byte(result), &encJson)

	return encJson, err
}

/**
	DEPRECATED
 */
func SetCacheStruct(namespace string, key []byte, value interface{}) error {
	if Client == nil {
		return errors.New("cache not connet")
	}
	nameSpaceKey := hex.EncodeToString(key)

	jsonM, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return Client.Set(namespace+":"+nameSpaceKey, jsonM, 0).Err()
}

/**
	DEPRECATED
 */
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

/**
	DEPRECATED
 */
func SetSafeStruct(namespace string, key []byte, value interface{}) error {
	nameSpaceKey := hex.EncodeToString(key)

	jsonM, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(namespace+":"+nameSpaceKey, jsonM, 0).Err()
}