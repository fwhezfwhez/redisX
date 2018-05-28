package redisX

import "log"
//1. more details referring '1728565484@qq.com' or submit your issue on github

//2. redisX is based on https://github.com/garyburd/redigo/

//3. redisX supports pool and transaction options

//4. the transaction object is service oriented ,which means you don't have to care whether it is open,thus no need to
//make a Close() and Begin() for this transaction.We only use Commit() and RollBack()  to guarantee several operations
//to be atomic.By the way,you can call Begin() and Close() to make a transaction complete.

//5. resultParser use origin https://github.com/garyburd/redigo/redis
func init(){
	log.SetFlags(log.LstdFlags | log.Llongfile)
}