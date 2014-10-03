package api

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Helper function that normalizes structs for comparison with reflect.DeepEqual
func normalize(v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Could not normalize object %+v due to JSON marshalling error: %s", v, err))
	}
	err = json.Unmarshal(j, v)
	if err != nil {
		panic(fmt.Sprintf("Could not normalize object %+v due to JSON un-marshalling error: %s", v, err))
	}
}

func normalizeDeepEqual(u, v interface{}) bool {
	normalize(u)
	normalize(v)
	return reflect.DeepEqual(u, v)
}
