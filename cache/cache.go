package cache

type Cache interface {
	Init(config interface{}) (error)
	Get(key interface{}) (val interface{})
	Set(key interface{}, data... interface{}) (bool, error)
	Delete(key interface{}) (bool, error)
	Flush() (bool, error)
}