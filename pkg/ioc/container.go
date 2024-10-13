package ioc

var _ Container = &MapContainer{
	name:   "controller",
	storge: make(map[string]Object),
}
