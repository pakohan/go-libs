package flatscan

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

type Entity interface {
	Key() string
	Type() string
	AEKey(c appengine.Context) *datastore.Key
}

type Counter struct {
	ZipCodes int
	Offers   int
}

func (c Counter) Key() string {
	return "counter"
}

func (c Counter) Type() string {
	return "counter"
}

func (c Counter) AEKey(con appengine.Context) *datastore.Key {
	return datastore.NewKey(con, "counter", "counter", 0, nil)
}

func InitCounter(c appengine.Context) {
	amount, _ := datastore.NewQuery("FlatOffer").Count(c)
	amountZipcodes, _ := datastore.NewQuery("ZIP").Count(c)

	counter := Counter{amountZipcodes, amount}
	datastore.Put(c, counter.AEKey(c), counter)
}

func StoreEntity(c appengine.Context, e Entity) {
	item := &memcache.Item{Key: e.Key(), Object: e}

	memcache.Add(c, item)

	datastore.Put(c, e.AEKey(c), e)
}

func Exists(c appengine.Context, e Entity) bool {
	item := getFromCache(c, e)
	if item != nil {
		return true
	}

	amount, err := datastore.NewQuery(e.Type()).Filter("__key__ =", e.AEKey(c)).Count(c)
	if amount > 0 && err == nil {
		return true
	}

	return false
}

func getFromCache(c appengine.Context, e Entity) interface{} {
	item, err := memcache.Get(c, e.Key())
	if err == nil && item != nil {
		return item.Object
	}

	return nil
}
