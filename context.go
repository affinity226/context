// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"sync"
)

var (
	mutex sync.RWMutex
	data  = make(map[string]map[string]interface{})
)

// Set stores a value for a given key in a given request.
func Set(idx, key string, val interface{}) {
	mutex.Lock()
	if data[idx] == nil {
		data[idx] = make(map[string]interface{})
	}
	data[idx][key] = val
	mutex.Unlock()
}

// Get returns a value stored for a given key in a given request.
func Get(idx, key string) interface{} {
	mutex.RLock()
	if ctx := data[idx]; ctx != nil {
		value, ok := ctx[key]
		mutex.RUnlock()
		if ok {
			return value
		} else {
			return nil
		}
	}
	mutex.RUnlock()
	return nil
}

// GetOk returns stored value and presence state like multi-value return of map access.
func GetOk(idx, key string) (interface{}, bool) {
	mutex.RLock()
	if _, ok := data[idx]; ok {
		value, ok := data[idx][key]
		mutex.RUnlock()
		return value, ok
	}
	mutex.RUnlock()
	return nil, false
}

// GetAll returns all stored values for the request as a map. Nil is returned for invalid requests.
func GetAll(idx string) map[string]interface{} {
	mutex.RLock()
	if context, ok := data[idx]; ok {
		result := make(map[string]interface{}, len(context))
		for k, v := range context {
			result[k] = v
		}
		mutex.RUnlock()
		return result
	}
	mutex.RUnlock()
	return nil
}

// GetAllOk returns all stored values for the request as a map and a boolean value that indicates if
// the request was registered.
func GetAllOk(idx string) (map[string]interface{}, bool) {
	mutex.RLock()
	context, ok := data[idx]
	result := make(map[string]interface{}, len(context))
	for k, v := range context {
		result[k] = v
	}
	mutex.RUnlock()
	return result, ok
}

// Delete removes a value stored for a given key in a given request.
func Delete(idx, key string) {
	mutex.Lock()
	if data[idx] != nil {
		delete(data[idx], key)
	}
	mutex.Unlock()
}

// Clear removes all values stored for a given request.
//
// This is usually called by a handler wrapper to clean up request
// variables at the end of a request lifetime. See ClearHandler().
func Clear(idx string) {
	mutex.Lock()
	clear(idx)
	mutex.Unlock()
}

// clear is Clear without the lock.
func clear(idx string) {
	delete(data, idx)
}

func PrintData(idx string) {
	mutex.RLock()
	if data, ok := data[idx]; ok {
		for key, value := range data {
			fmt.Printf("%s %s %v\n", idx, key, value)
		}
	}
	mutex.RUnlock()
}
func PrintAll() {
	mutex.RLock()
	for idx, data := range data {
		for key, value := range data {
			fmt.Printf("%s %s %v\n", idx, key, value)
		}
	}
	mutex.RUnlock()
}
