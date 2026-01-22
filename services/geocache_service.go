package services

import "sync"

var geoCache = struct {
	sync.RWMutex
	data map[string]Coordinates
}{
	data: make(map[string]Coordinates),
}

func GetCachedCoordinates(address string) (Coordinates, bool) {
	geoCache.RLock()
	defer geoCache.RUnlock()
	coords, ok := geoCache.data[address]
	return coords, ok
}

func SetCachedCoordinates(address string, coords Coordinates) {
	geoCache.Lock()
	defer geoCache.Unlock()
	geoCache.data[address] = coords
}
