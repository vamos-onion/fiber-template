package interfaces

type RedisQuery interface {
	HSET(string, string, interface{}) error
	HGET(string, string) (string, error)
	HDEL(string, string) error
	HKEYS(string) []string
}

func HSET(p RedisQuery, m string, n string, i interface{}) error {
	return p.HSET(m, n, i)
}

func HGET(p RedisQuery, m string, n string) (string, error) {
	return p.HGET(m, n)
}

func HDEL(p RedisQuery, m string, n string) error {
	return p.HDEL(m, n)
}

func HKEYS(p RedisQuery, m string) []string {
	return p.HKEYS(m)
}
