package product

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"helper/redisCli"
)

func BestSelling() *Sales {
	return &Sales{Tag: fmt.Sprint("tag:", "bestselling:"), ProductKey: "is_hot"}
}

func NewArrivals() *Sales {
	return &Sales{Tag: fmt.Sprint("tag:", "newarrivals:"), ProductKey: "is_new"}
}

type Sales struct {
	Tag        string
	ProductKey string
}

func GetItemSaleTags(item_id int) map[string]int {
	tags := map[string]int{}
	r := redisCli.Conn()
	defer r.Close()
	s, _ := redis.Ints(r.Do("HMGET", fmt.Sprint("product:", item_id), "is_hot", "is_new"))
	tags["is_hot"] = s[0]
	tags["is_new"] = s[1]
	return tags
}

func (sale *Sales) Add(item_id, sort int) {
	r := redisCli.Conn()
	defer r.Close()
	r.Do("ZADD", sale.Tag, sort, item_id)
	r.Do("HSET", fmt.Sprint("product:", item_id), sale.ProductKey, 1)
}

func (sale *Sales) Rem(item_ids ...int) {
	r := redisCli.Conn()
	defer r.Close()
	for _, item_id := range item_ids {
		r.Do("ZREM", sale.Tag, item_id)
		r.Do("HDEL", fmt.Sprint("product:", item_id), sale.ProductKey)
	}
}

func (sale *Sales) Listing(offset, rowcount int) ([]*Item, int) {
	r := redisCli.Conn()
	defer r.Close()
	total := 0
	items := []*Item{}
	ids, err := redis.Ints(r.Do("ZRANGE", sale.Tag, offset, rowcount+offset))
	if err == nil {
		total, _ = redis.Int(r.Do("ZCARD", sale.Tag))
		for _, id := range ids {
			if item, err := get(fmt.Sprintf("id=%d",id)); err == nil {
				items = append(items, item)
			}
		}
		return items, total
	} else {
		fmt.Println("get sales listing failing:", err.Error(), sale.Tag)
		return items, total
	}
}
