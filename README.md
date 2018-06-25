**note**

1. more details referring '1728565484@qq.com' or submit your issue on github

2. redisX is based on https://github.com/garyburd/redigo/

3. redisX supports pool and transaction options

4. the transaction object is service oriented ,which means you don't have to care whether it is open,thus no need to
make a Close() and Begin() for this transaction.We only use Commit() and RollBack()  to guarantee several operations
to be atomic.By the way,you can call Begin() and Close() to make a transaction complete.

5. resultParser use origin https://github.com/garyburd/redigo/redis

6. transaction for redisx  serves for redis server itself,when do SET HSET DEL HDEL in a transaction,this transaction
only makes sure no data actually saved in redis server,however data in a GET,HGET command has been set into your program.
For examle:
 ```go
    data := redis.String(tran.Do("Get","key"))
    tran.Do("Set","key2","value")
    tran.RollBack()
 ```
 **data value has been set ,but key2 has been rollBack.**

**start**

go get github.com/fwhezfwhez/redisX

**example**
```go
package main

import (
	"github.com/fwhezfwhez/redisX"
	"log"
)

func init(){
	log.SetFlags(log.LstdFlags | log.Llongfile)
}
func main(){
	//
	redisDb:= redisX.RedisX{}
	redisDb.DataSource("redis://localhost:6379")



	tran := redisDb.BeginTran()
	tran.Begin()


	defer tran.Close()
	_,er:=tran.Do("SET","username","ft2")
	if er!=nil{
		log.Println(er.Error())
		tran.RollBack()
		return
	}
	_,er=tran.Do("HSET","user","username2","ft2")

	if er!=nil{
		log.Println(er.Error())
		tran.RollBack()
		return
	}

	_,er=tran.Do("DEL","username")

	if er!=nil{
		log.Println(er.Error())
		tran.RollBack()
		return
	}
	//tran.RollBack()
	_,er=tran.Do("HDEL","user","username")

	if er!=nil{
		log.Println(er.Error())
		tran.RollBack()
		return
	}
	//tran.RollBack()
	rs,er:=tran.Do("HGET","user","username2")
	if er!=nil{
		log.Println(er.Error())
		tran.RollBack()
		return
	}
	tran.Commit()
	log.Println(redisX.String(rs,nil))
}

```