package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/pubsub"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mediocregopher/radix.v2/sentinel"
	"github.com/mediocregopher/radix.v2/util"
)

var sentinelPool *sentinel.Client
var securitylPool *sentinel.Client
var redisPool *pool.Pool

var subChannelName string
var subChannelNameAgent string

func InitiateRedis() {

	var err error

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}
		if err = client.Cmd("AUTH", redisPassword).Err; err != nil {
			client.Close()
			return nil, err
		}
		if err = client.Cmd("select", redisDb).Err; err != nil {
			client.Close()
			return nil, err
		}
		return client, nil
	}

	if redisMode == "sentinel" {
		sentinelIps := strings.Split(sentinelHosts, ",")

		if len(sentinelIps) > 1 {
			sentinelIp := fmt.Sprintf("%s:%s", sentinelIps[0], sentinelPort)
			sentinelPool, err = sentinel.NewClientCustom("tcp", sentinelIp, 10, df, redisClusterName)

			if err != nil {
				errHndlrNew("InitiateRedis", "InitiateSentinel", err)
			}

			securitylPool, err = sentinel.NewClientCustom("tcp", sentinelIp, 10, df, redisClusterName)

			if err != nil {
				errHndlrNew("InitiateRedis", "InitiateSentinel-securitylPool", err)
			}
		} else {
			fmt.Println("Not enough sentinel servers")
		}
	} else {
		redisPool, err = pool.NewCustom("tcp", redisIp, 10, df)

		if err != nil {
			errHndlrNew("InitiateRedis", "InitiatePool", err)
		}
	}
}

// Redis String Methods
func RedisAdd(key, value string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	isExists, _ := client.Cmd("EXISTS", key).Int()

	if isExists == 1 {
		return "Key Already exists"
	} else {
		result, sErr := client.Cmd("set", key, value).Str()
		errHndlr(sErr)
		DialerLog(result)
		return result
	}
}

func RedisSet(key, value string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("set", key, value).Str()
	errHndlr(sErr)
	DialerLog(result)
	return result
}

func RedisSetNx(key, value string) int {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("setnx", key, value).Int()
	errHndlr(sErr)
	DialerLog(fmt.Sprintf("%d", result))
	return result
}

func RedisGet(key string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	strObj, _ := client.Cmd("get", key).Str()
	DialerLog(fmt.Sprintf("%+v", strObj))
	return strObj
}

func AppendIfMissing(windowList []string, i string) []string {
	for _, ele := range windowList {
		if ele == i {
			return windowList
		}
	}
	return append(windowList, i)
}

func RedisSearchKeys(pattern string) []string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSearchKeys", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	matchingKeys := make([]string, 0)

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	DialerLog(fmt.Sprintf("Start ScanAndGetKeys:: %s", pattern))
	scanResult := util.NewScanner(client, util.ScanOpts{Command: "SCAN", Pattern: pattern, Count: 1000})

	for scanResult.HasNext() {
		//fmt.Println("next:", scanResult.Next())
		matchingKeys = AppendIfMissing(matchingKeys, scanResult.Next())
	}

	DialerLog(fmt.Sprintf("Scan Result:: %+v", matchingKeys))
	return matchingKeys
}

func RedisIncr(key string) int {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("incr", key).Int()
	errHndlr(sErr)
	DialerLog(fmt.Sprintf("%d", result))
	return result
}

func RedisIncrBy(key string, value int) int {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisSet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, sErr := client.Cmd("incrby", key, value).Int()
	errHndlr(sErr)
	DialerLog(fmt.Sprintf("%d", result))
	return result
}

func RedisRemove(key string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisRemove", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tempResult, sErr := client.Cmd("del", key).Int()

	errHndlr(sErr)
	DialerLog(fmt.Sprintf("%d", tempResult))
	if tempResult == 1 {
		return true
	} else {
		return false
	}
}

func RedisCheckKeyExist(key string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in CheckKeyExist", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tempResult, sErr := client.Cmd("exists", key).Int()
	errHndlr(sErr)
	DialerLog(fmt.Sprintf("%d", tempResult))
	if tempResult == 1 {
		return true
	} else {
		return false
	}
}

// Redis Hashes Methods

func RedisHashGetAll(hkey string) map[string]string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetAll", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}
	strHash, _ := client.Cmd("hgetall", hkey).Map()
	bytes, err := json.Marshal(strHash)
	if err != nil {
		fmt.Println(err)
	}
	text := string(bytes)
	DialerLog(text)
	return strHash
}

func RedisHashGetField(hkey, field string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashGetAll", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}
	strValue, _ := client.Cmd("hget", hkey, field).Str()
	if err != nil {
		fmt.Println(err)
	}
	DialerLog(strValue)
	return strValue
}

func RedisHashSetField(hkey, field, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashSetField", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	tempResult, _ := client.Cmd("hset", hkey, field, value).Int()

	if tempResult == 1 {
		return true
	} else {
		return false
	}
}

func RedisHashSetMultipleField(hkey string, data map[string]string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisHashSetField", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}
	DialerLog(fmt.Sprintf("%+v", data))
	for key, value := range data {
		client.Cmd("hset", hkey, key, value)
	}
	//fmt.Println(true)
	return true
}

// Redis List Methods

func RedisListLpop(lname string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	lpopItem, _ := client.Cmd("lpop", lname).Str()
	DialerLog(lpopItem)
	return lpopItem
}

func RedisListLpush(lname, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpush", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, _ := client.Cmd("lpush", lname, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func RedisListRpush(lname, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpush", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, _ := client.Cmd("rpush", lname, value).Int()
	if result > 0 {
		return true
	} else {
		return false
	}
}

func RedisListLlen(lname string) int {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLlen", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	result, _ := client.Cmd("llen", lname).Int()
	return result
}

func SecurityGet(key string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				securitylPool.PutMaster(redisClusterName, client)
			}
			//			else {
			//				redisPool.Put(client)
			//			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = securitylPool.GetMaster(redisClusterName)
		client.Cmd("select", "0")
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer securitylPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redis.DialTimeout("tcp", securityIp, time.Duration(10)*time.Second)
		errHndlr(err)
		//defer client.Close()

		//authServer
		authE := client.Cmd("auth", redisPassword)
		errHndlr(authE.Err)
	}

	strObj, _ := client.Cmd("get", key).Str()
	//fmt.Println(strObj)
	return strObj
}

// Redis Set Methods

func RedisSetAdd(key string, members []string) string {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	saddResult, _ := client.Cmd("sadd", key, members).Str()
	DialerLog(fmt.Sprintf("saddResult : %s :: %s", key, saddResult))
	return saddResult
}

func RedisSetIsMember(key, value string) bool {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLpop", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	sismemberResult, _ := client.Cmd("sismember", key, value).Int()
	DialerLog(fmt.Sprintf("sismember : %s:%s :: %s", key, value, sismemberResult))

	if sismemberResult == 1 {
		return true
	} else {
		return false
	}
}

// Redis PubSub

func PubSub() {
	var c2 *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PubSub", r)
		}

		if c2 != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, c2)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	subChannelName = fmt.Sprintf("dialer%s", dialerId)
	for {
		if redisMode == "sentinel" {

			c2, err = sentinelPool.GetMaster(redisClusterName)
			errHndlrNew("PubSub", "getConnFromPool", err)
			//defer sentinelPool.PutMaster(redisClusterName, c2)

			psc := pubsub.NewSubClient(c2)
			psr := psc.Subscribe(subChannelName)
			//ppsr := psc.PSubscribe("*")

			//if ppsr.Err == nil {

			for {
				psr = psc.Receive()
				if psr.Timeout() {
					fmt.Println("psc.Receive Timeout:: ", psr.Timeout())
					break

				}
				if psr.Err != nil {

					fmt.Println("psc.Receive Err:: ", psr.Err.Error())
					break
				}

				var subEvent = SubEvents{}
				json.Unmarshal([]byte(psr.Message), &subEvent)
				go OnEvent(subEvent)
			}
			//s := strings.Split("127.0.0.1:5432", ":")
			//}

			psc.Unsubscribe(subChannelName)

		} else {
			c2, err = redis.Dial("tcp", redisIp)
			errHndlr(err)
			defer c2.Close()

			//authServer
			authE := c2.Cmd("auth", redisPassword)
			errHndlr(authE.Err)

			psc := pubsub.NewSubClient(c2)
			psr := psc.Subscribe(subChannelName)
			//ppsr := psc.PSubscribe("*")

			//if ppsr.Err == nil {

			for {
				psr = psc.Receive()
				if psr.Timeout() {
					fmt.Println("psc.Receive Timeout:: ", psr.Timeout())
					break

				}

				if psr.Err != nil {

					fmt.Println("psc.Receive Err:: ", psr.Err.Error())
					break
				}

				subEvent := SubEvents{}
				json.Unmarshal([]byte(psr.Message), &subEvent)
				go OnEvent(subEvent)
			}

			psc.Unsubscribe(subChannelName)

		}
		time.Sleep(1 * time.Second)
	}
}

func PubSubAgentChan() {
	var c3 *redis.Client
	var err1 error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in PubSub", r)
		}

		if c3 != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, c3)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	subChannelNameAgent = fmt.Sprintf("dialerAgent%s", dialerId)
	for {
		if redisMode == "sentinel" {

			c3, err1 = sentinelPool.GetMaster(redisClusterName)
			errHndlrNew("PubSub", "getConnFromPool", err1)
			//defer sentinelPool.PutMaster(redisClusterName, c2)

			psc := pubsub.NewSubClient(c3)
			psr := psc.Subscribe(subChannelNameAgent)
			//ppsr := psc.PSubscribe("*")

			//if ppsr.Err == nil {

			for {
				psr = psc.Receive()
				if psr.Timeout() {
					fmt.Println("psc.Receive Timeout:: ", psr.Timeout())
					break

				}
				if psr.Err != nil {

					fmt.Println("psc.Receive Err:: ", psr.Err.Error())
					break
				}

				subEvent := SubEvents{}
				json.Unmarshal([]byte(psr.Message), &subEvent)
				go OnEventAgent(subEvent)
			}
			//s := strings.Split("127.0.0.1:5432", ":")
			//}

			psc.Unsubscribe(subChannelNameAgent)

		} else {
			c3, err1 = redis.Dial("tcp", redisIp)
			errHndlr(err1)
			defer c3.Close()

			//authServer
			authE := c3.Cmd("auth", redisPassword)
			errHndlr(authE.Err)

			psc := pubsub.NewSubClient(c3)
			psr := psc.Subscribe(subChannelNameAgent)
			//ppsr := psc.PSubscribe("*")

			//if ppsr.Err == nil {

			for {
				psr = psc.Receive()
				if psr.Timeout() {
					fmt.Println("psc.Receive Timeout:: ", psr.Timeout())
					break

				}

				if psr.Err != nil {

					fmt.Println("psc.Receive Err:: ", psr.Err.Error())
					break
				}

				var subEvent = SubEvents{}
				json.Unmarshal([]byte(psr.Message), &subEvent)
				go OnEventAgent(subEvent)
			}

			psc.Unsubscribe(subChannelNameAgent)

		}
		time.Sleep(1 * time.Second)
	}
}

func Publish(channel, message string) {
	var client *redis.Client
	var err error

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisListLlen", r)
		}

		if client != nil {
			if redisMode == "sentinel" {
				sentinelPool.PutMaster(redisClusterName, client)
			} else {
				redisPool.Put(client)
			}
		} else {
			fmt.Println("Cannot Put invalid connection")
		}
	}()

	if redisMode == "sentinel" {
		client, err = sentinelPool.GetMaster(redisClusterName)
		errHndlrNew("OnEvent", "getConnFromSentinel", err)
		//defer sentinelPool.PutMaster(redisClusterName, client)
	} else {
		client, err = redisPool.Get()
		errHndlrNew("OnEvent", "getConnFromPool", err)
		//defer redisPool.Put(client)
	}

	client.Cmd("publish", channel, message)
}
