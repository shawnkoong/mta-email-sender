package main

import (
	"sync"
	"time"
)

type EmailTracker struct {
	mux      sync.RWMutex
	emailMap map[string]*RouteTracker
}

type RouteTracker struct {
	mux      sync.RWMutex
	routeMap map[string]time.Time
	// ideally have routeMap as map[string]map[string]timestamp as route to alert to timestamp
	// but then would periodically have to clean the map of alerts to not hold all old alerts
}

var instance *EmailTracker
var once sync.Once

func getEmailTracker() *EmailTracker {
	once.Do(func() {
		instance = &EmailTracker{
			emailMap: make(map[string]*RouteTracker),
		}
	})
	return instance
}

func NewRouteTracker() *RouteTracker {
	return &RouteTracker{
		routeMap: make(map[string]time.Time),
	}
}

func (et *EmailTracker) update(email string, routeTracker *RouteTracker) {
	et.mux.Lock()
	defer et.mux.Unlock()
	et.emailMap[email] = routeTracker
}

func (et *EmailTracker) get(email string) (*RouteTracker, bool) {
	et.mux.RLock()
	defer et.mux.RUnlock()
	rt, ok := et.emailMap[email]
	return rt, ok
}

func (rt *RouteTracker) update(route string, now time.Time) {
	rt.mux.Lock()
	defer rt.mux.Unlock()
	rt.routeMap[route] = now
}

func (rt *RouteTracker) get(route string) (time.Time, bool) {
	rt.mux.RLock()
	defer rt.mux.RUnlock()
	t, ok := rt.routeMap[route]
	return t, ok
}

// returns true if new email should be sent, otherwise returns false
func (rt *RouteTracker) checkLastTimeSent(route string) bool {
	lastTime, ok := rt.get(route)
	if !ok {
		return true
	}
	timeDiff := time.Since(lastTime)
	if timeDiff > 15*time.Minute { // threshold of 15 minutes for now
		return true
	}
	return false
}
