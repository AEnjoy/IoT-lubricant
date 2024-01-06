package ioc

var Controller Container = &MapContainer{
	name:   "controller",
	storge: make(map[string]Object),
}

var Api Container = &MapContainer{
	name:   "api",
	storge: make(map[string]Object),
}
