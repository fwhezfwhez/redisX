package redisX

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
	"errors"
)


type RedisX struct {
	Pool     *redis.Pool
	LocalCon Connection
}

type Connection struct {
	Conn redis.Conn
}
type Transaction struct {
	Con    Connection
	LastKV map[string]interface{} //command-string      kv interface{}
}

//init a redis datasource by default config
//if you want to config more use api RedisX.Config(time.Duration,int,int)
//if you want to config more specifically use it like
/*
	 redisx :=RedisX{}
	 redisx.DataSource("redis://localhost:6379")
	 redisx.pool.MaxIdle=xxx
	 redis.pool.XX=xxx
*/
func (r *RedisX) DataSource(dataSource string) {
	//"redis://localhost:6379"
	r.Pool = &redis.Pool{
		MaxIdle: 200,
		MaxActive:  200,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(dataSource)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	//&{0xc042092080 0xc0420ca000 0}
	r.LocalCon = Connection{r.Pool.Get()}
}

//config
func (r RedisX) Config(idleTimeout time.Duration, maxIdle, maxActive int) {
	r.Pool.MaxIdle = maxIdle
	r.Pool.IdleTimeout = idleTimeout
	r.Pool.MaxActive = maxActive
}

//begin a transaction
func (r RedisX) BeginTran() *Transaction {
	lastKV := make(map[string]interface{})
	SET:=make(map[interface{}]interface{})
	HSET:=make(map[interface{}]interface{})
	lastKV["HSET"]=HSET
	lastKV["SET"]=SET
	fmt.Println(r.Pool)
	con :=Connection{Conn: r.Pool.Get()}
	return &Transaction{
		Con:con,
		LastKV: lastKV,
	}
}


//do a command
func (c Connection) Do(command string, args ...interface{}) (interface{}, error) {
	return c.Conn.Do(command, args...)
}

//transaction begins
func (tran *Transaction) Begin() {
	fmt.Println("begin a redis transaction")
}

//do a command in a transaction  which can rollback when faced with an error
func (tran *Transaction) Do(command string, args ...interface{}) (interface{}, error) {
	switch command{
	case "SET","DEL":
		setMap :=(tran.LastKV["SET"]).(map[interface{}]interface{})
		if checkLength(command,args...){
			temp,er :=tran.Con.Do("GET",args[0])
			if er!=nil{
				    setMap[args[0]]=""
			}
			setMap[args[0]]=temp
		}else{
			return nil,errors.New(command+"参数数量错误")
		}
	case "HSET","HDEL":
		if checkLength(command,args...) {
			hSetMap := (tran.LastKV["HSET"]).(map[interface{}]interface{})
			temp, _ := tran.Con.Do("HGET", args[0], args[1])
			var inMap map[interface{}]interface{}
			if hSetMap[args[0]] != nil {
				inMap = hSetMap[args[0]].(map[interface{}]interface{})
			} else {
				inMap = make(map[interface{}]interface{})
			}
			inMap[args[1]] = temp
			hSetMap[args[0]] = inMap
		}else{
			return nil,errors.New(command+"参数数量错误")
		}
	//default:
	//	return  tran.Con.Do(command, args...)
	}
	return tran.Con.Do(command, args...)
}

//to make cached lastKV flushed
func (tran *Transaction) Commit() {
	tran.LastKV = make(map[string]interface{})
}

//rollBack
func (tran *Transaction) RollBack() {
	if len(tran.LastKV)!=0{
		setMap :=(tran.LastKV["SET"]).(map[interface{}]interface{})
		if len(setMap)!=0{
			for k,v:=range setMap{
				tran.Con.Do("SET",k,v)
			}
		}
		hSetMap := (tran.LastKV["HSET"]).(map[interface{}]interface{})
		if len(hSetMap)!=0{
			for mainKey,v:=range hSetMap{
				inMap:=v.(map[interface{}]interface{})
				for subKey,value:=range inMap{
					tran.Con.Do("HSET",mainKey,subKey,value)
				}
			}
		}
	}
}

//
func (tran *Transaction) Close() {
	fmt.Println("end a redis transaction")
}


//only check for SET HSET DEL
func checkLength(command string,args...interface{}) bool{
	length :=len(args)
	switch command {
	case "SET":
		return length==2
	case "HSET":
		return length==3
	case "DEL":
		return length==1
	case "HDEL":
		return length==2
	}
	return true
}


//Abandoned
func GetRedis(url string) *redis.Pool {
	fmt.Println("get redis url",url)
	var rs *redis.Pool
	rs= &redis.Pool{
		MaxIdle: 200,
		//MaxActive:   0,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(url)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Println(rs)
	return rs
}

//&{0x622690 0x6228f0 200 0 3m0s false {0 0} <nil> false 0 {{<nil> <nil> <nil> <nil>} 0}}