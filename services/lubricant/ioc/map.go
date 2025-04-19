package ioc

import (
	"fmt"
	"sort"
	"sync"

	"github.com/aenjoy/iot-lubricant/pkg/logger"
)

var _ Container = (*MapContainer)(nil)

type MapContainer struct {
	name    string
	storge  map[string]Object
	l       sync.Mutex
	inited  bool
	showLog bool
}

func (c *MapContainer) Version(name string) string {
	if obj, ok := c.storge[name]; ok {
		return obj.Version()
	}
	return ""
}
func (c *MapContainer) LoadObject(s map[string]Object) {
	logger.Debugf("Load objects: %v", s)
	c.l.Lock()
	defer c.l.Unlock()
	c.storge = s
}
func (c *MapContainer) Registry(name string, obj Object) {
	logger.Debugf("Registry object: %s", name)
	c.l.Lock()
	defer c.l.Unlock()
	c.storge[name] = obj
}

func (c *MapContainer) Get(name string) any {
	return c.storge[name]
}

func (c *MapContainer) Init() error {
	if c.inited {
		return fmt.Errorf("object %s all has already been initialized", c.name)
	}
	// Create a slice of object names and weights for sorting
	type weightedObject struct {
		name   string
		object Object
	}

	var weightedObjects []weightedObject

	for name, obj := range c.storge {
		weightedObjects = append(weightedObjects, weightedObject{name: name, object: obj})
	}

	// Sort by weight (ascending)
	sort.Slice(weightedObjects, func(i, j int) bool {
		return weightedObjects[i].object.Weight() < weightedObjects[j].object.Weight()
	})

	// Initialize objects in the order of their weights
	for _, wObj := range weightedObjects {
		logger.Debugf("[%s] %s init start with weight %d", c.name, wObj.name, wObj.object.Weight())
		if err := wObj.object.Init(); err != nil {
			return fmt.Errorf("%s init error, %s", wObj.name, err)
		}
		if c.showLog {
			logger.Infof("[%s] %s init success with weight %d", c.name, wObj.name, wObj.object.Weight())
		}
	}

	c.inited = true
	return nil
}
