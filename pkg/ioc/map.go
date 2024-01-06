package ioc

import (
	"fmt"
	"sort"
)

var _ Container = (*MapContainer)(nil)

type MapContainer struct {
	name   string
	storge map[string]Object
}

func (c *MapContainer) Registry(name string, obj Object) {
	c.storge[name] = obj
}

func (c *MapContainer) Get(name string) any {
	return c.storge[name]
}

func (c *MapContainer) Init() error {
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
		if err := wObj.object.Init(); err != nil {
			return fmt.Errorf("%s init error, %s", wObj.name, err)
		}
		fmt.Printf("[%s] %s init success with weight %d\n", c.name, wObj.name, wObj.object.Weight())
	}

	return nil
}
