package helpers

import (
	"maps"
	"slices"
	"sync"
)

type SafeMap[keyType comparable, valueType any] struct {
	data  map[keyType]valueType
	mutex *sync.RWMutex
}

func NewSafeMap[keyType comparable, valueType any](data map[keyType]valueType) *SafeMap[keyType, valueType] {
	return &SafeMap[keyType, valueType]{
		data:  data,
		mutex: &sync.RWMutex{},
	}
}

func (safeMap *SafeMap[keyType, valueType]) Get(key keyType) (valueType, bool) {
	defer safeMap.mutex.RUnlock()
	safeMap.mutex.RLock()

	value, ok := safeMap.data[key]

	return value, ok
}

func (safeMap *SafeMap[keyType, valueType]) Set(key keyType, value valueType) {
	defer safeMap.mutex.Unlock()
	safeMap.mutex.Lock()

	safeMap.data[key] = value
}

func (safeMap *SafeMap[keyType, valueType]) GetKeys() []keyType {
	defer safeMap.mutex.RUnlock()
	safeMap.mutex.RLock()

	return slices.Collect(maps.Keys(safeMap.data))
}

func (safeMap *SafeMap[keyType, valueType]) GetValues() []valueType {
	defer safeMap.mutex.RUnlock()
	safeMap.mutex.RLock()

	return slices.Collect(maps.Values(safeMap.data))
}
