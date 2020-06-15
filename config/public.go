package config

const (
	CachePath = "cache.db"
)

func SetEmail(email *Email) error {
	cache, err := NewCache(CachePath)
	if err != nil {
		return err
	}
	defer cache.Close()

	return cache.SetEmail(email)
}

func SetNeo(neo *Neo) error {
	cache, err := NewCache(CachePath)
	if err != nil {
		return err
	}
	defer cache.Close()

	return cache.SetNeo(neo)
}
