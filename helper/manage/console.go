package manage

import "helper/redisCli"

type Console int

func (*Console)RedisActiveCount()  int {
	return redisCli.ActiveCount()
}