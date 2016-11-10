// +build darwin
// +build arm arm64

package app

import (
	"runtime"
	"sync"
	"time"
)

func init() {
	// Lock the goroutine responsible for initialization to an OS thread.
	// This means the goroutine running main (and calling the run function
	// below) is locked to the OS thread that started the program. This is
	// necessary for the correct delivery of UIKit events to the process.
	//
	// A discussion on this topic:
	// https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()
}

type Signal struct{}

type AppDelegate interface {
	InitDone()
	LifecycleEvents() <-chan LifecycleEvent
	VSync() <-chan Signal

	HandleConfigurationEvents(out chan<- ConfigurationEvent)
	HandleTouchEvents(out chan<- TouchEvent)
}

var defaultApp = &appDelegate{
	lifecycleEvents: make(chan LifecycleEvent),
	vsyncEvents:     make(chan Signal),
	maxDispatchTime: 1 * time.Second,

	initWG: new(sync.WaitGroup),
	mux:    new(sync.RWMutex),
}

type appDelegate struct {
	// lifecycleEvents must be handled in real-time.
	lifecycleEvents chan LifecycleEvent
	// vsyncEvents must be handled in real-time.
	vsyncEvents chan Signal

	// maxDispatchTime sets the maximum time the send operation
	// allowed to wait while channel is blocked.
	maxDispatchTime time.Duration
	// channels below are optional and will be sent to only
	// if handled by an external client.

	configurationEvents chan<- ConfigurationEvent
	touchEvents         chan<- TouchEvent

	initWG *sync.WaitGroup
	mux    *sync.RWMutex
}

func (a *appDelegate) InitDone() {
	a.initWG.Done()
}

func (a *appDelegate) LifecycleEvents() <-chan LifecycleEvent {
	return a.lifecycleEvents
}

func (a *appDelegate) VSync() <-chan Signal {
	return a.vsyncEvents
}

func (a *appDelegate) HandleConfigurationEvents(out chan<- ConfigurationEvent) {
	a.mux.Lock()
	a.configurationEvents = out
	a.mux.Unlock()
}

func (a *appDelegate) getConfigurationEventsOut() chan<- ConfigurationEvent {
	a.mux.RLock()
	out := a.configurationEvents
	a.mux.RUnlock()
	return out
}

func (a *appDelegate) HandleTouchEvents(out chan<- TouchEvent) {
	a.mux.Lock()
	a.touchEvents = out
	a.mux.Unlock()
}

func (a *appDelegate) getTouchEventsOut() chan<- TouchEvent {
	a.mux.RLock()
	out := a.touchEvents
	a.mux.RUnlock()
	return out
}
