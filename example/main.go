package main

import (
	"redisX"
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
