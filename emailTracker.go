package main

import "sync"

type EmailTracker struct {
	mtx        sync.Mutex
	emailsSent map[string]string
}

var instance *EmailTracker
var once sync.Once

func getEmailTracker() *EmailTracker {
	once.Do(func() {
		instance = &EmailTracker{
			emailsSent: make(map[string]string),
		}
	})
	return instance
}

func (tracker *EmailTracker) update() {
	tracker.mtx.Lock()
	defer tracker.mtx.Unlock()
	// update logic
}
