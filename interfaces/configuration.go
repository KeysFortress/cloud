package interfaces

type Configuration interface {
	Load() bool
	GetKey(key string) interface{}
}
