package apiserver

type Config struct {
	Addr      string `toml:"bind_addr"`
	DBURL     string `toml:"database_url"`
	JWTSecret string `toml:"jwt_secret"`
}
