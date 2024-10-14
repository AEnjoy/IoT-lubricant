package ioc

var Controller Container = &MapContainer{
	name:   "controller",
	storge: make(map[string]Object),
}
