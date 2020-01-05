package module

type Config struct {

}

type HTTPClientConfig struct {
	URL string `json:"url"`
	Domain string `json:"domain"`
	Cookie string `json:"cookie"`
}
