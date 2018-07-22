package webkit2

import (
	"github.com/ghetzel/go-stockutil/stringutil"
)

var registeredObjects = make(map[string]interface{})

func cgoregister(key string, obj interface{}) string {
	if key == `` || key == `auto` {
		key = stringutil.UUID().Base58()
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
