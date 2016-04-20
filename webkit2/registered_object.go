package webkit2

import (
	"github.com/satori/go.uuid"
)

var registeredObjects = make(map[string]interface{})

func cgoregister(key string, obj interface{}) string {
	if key == `` || key == `auto` {
		key = uuid.NewV4().String()
	}

	registeredObjects[key] = obj
	return key
}

func cgoget(key string) (interface{}, bool) {
	v, ok := registeredObjects[key]
	return v, ok
}

func cgounregister(key string) {
	delete(registeredObjects, key)
}
