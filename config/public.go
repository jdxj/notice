package config

const (
	CachePath = "cache.db"
)

var (
	DataStorage *Cache
)

func init() {
	var err error
	DataStorage, err = NewCache(CachePath)
	if err != nil {
		panic(err)
	}
}
