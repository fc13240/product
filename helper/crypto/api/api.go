package api

import (
	"helper/configs"
	"helper/crypto"
	"helper/net"
)

func init() {
	net.Do("/api.crypto.md5", func(act *net.Act) {
		act.Succ(configs.M{"to": crypto.Md5(act.Get("s"))})
	})
}
