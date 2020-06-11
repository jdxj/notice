package config

const (
	CachePath = "cache.db"
)

var (
	cache = NewCache(CachePath)
)

func Close() error {
	return cache.Close()
}
func GetEmail() (*Email, error) {
	return cache.GetEmail()
}

func SetEmail(email *Email) error {
	return cache.SetEmail(email)
}

func GetNeo() (*Neo, error) {
	return cache.GetNeo()
}

func SetNeo(neo *Neo) error {
	return cache.SetNeo(neo)
}
