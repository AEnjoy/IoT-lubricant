package ioc

type Container interface {
	Registry(name string, obj Object)
	Get(name string) any
	Init() error
}

type Object interface {
	Init() error
	Weight() uint16 // Add a Weight method to get the weight of the object
}
