package proxy

type Config struct {
	Host     string
	Port     int
	Username string
	Password string

	DialHost string
	DialPort int
}
