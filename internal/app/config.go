package app

type Config struct {
	port    int
	workers int
}

func NewConfig(port int, workers int) Config {
	return Config{
		port:    port,
		workers: workers,
	}
}

func (c Config) Port() int {
	return c.port
}

func (c Config) Workers() int {
	return c.workers
}
