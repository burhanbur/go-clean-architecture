package redis

type IRedisClient interface {
	RetrieveById(id string) (interface{}, error)
	SetValue(key string, value interface{}) error
	DelKey(key string) error
}
