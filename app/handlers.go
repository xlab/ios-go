package app

import "C"
import (
	"log"
	"time"
)

type LifecycleEvent struct {
	View uintptr
	Kind LifecycleEventKind
}

type LifecycleEventKind string

// See also https://developer.apple.com/reference/uikit/uiapplicationdelegate
const (
	ViewDidLoad LifecycleEventKind = "viewDidLoad"

	// WillResignActive is sent when the application is about to move
	// from active to inactive state. This can occur for certain types of
	// temporary interruptions (such as an incoming phone call or SMS message)
	// or when the user quits the application and it begins the transition to
	// the background state. Use this method to pause ongoing tasks, disable
	// timers, and throttle down OpenGL ES frame rates. Games should use this
	// method to pause the game.
	WillResignActive LifecycleEventKind = "applicationWillResignActive"
	// DidEnterBackground signals that you should release shared
	// resources, save user data, invalidate timers, and store enough
	// application state information to restore your application to its current
	// state in case it is terminated later. If your application supports
	// background execution, this method is called instead of
	// WillTerminate: when the user quits.
	DidEnterBackground LifecycleEventKind = "applicationDidEnterBackground"
	// WillEnterForeground is sent as part of the transition from the
	// background to the inactive state; here you can undo many of the changes
	// made on entering the background.
	WillEnterForeground LifecycleEventKind = "applicationWillEnterForeground"
	// DidBecomeActive should restart any tasks that were paused (or
	// not yet started) while the application was inactive. If the application
	// was previously in the background, optionally refresh the user interface.
	DidBecomeActive LifecycleEventKind = "applicationDidBecomeActive"
	// WillTerminate is sent when the application is about to
	// terminate. Save data if appropriate. See also
	// DidEnterBackground.
	WillTerminate LifecycleEventKind = "applicationWillTerminate"
)

//export onVSync
func onVSync() {
	defaultApp.initWG.Wait()

	select {
	case defaultApp.vsyncEvents <- Signal{}:
	default:
	}
}

//export onViewDidLoad
func onViewDidLoad(view uintptr) {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		View: view,
		Kind: ViewDidLoad,
	}
	defaultApp.lifecycleEvents <- event
}

//export onApplicationWillResignActive
func onApplicationWillResignActive() {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Kind: WillResignActive,
	}
	defaultApp.lifecycleEvents <- event
}

//export onApplicationDidEnterBackground
func onApplicationDidEnterBackground() {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Kind: DidEnterBackground,
	}
	defaultApp.lifecycleEvents <- event
}

//export onApplicationWillEnterForeground
func onApplicationWillEnterForeground() {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Kind: WillEnterForeground,
	}
	defaultApp.lifecycleEvents <- event
}

//export onApplicationDidBecomeActive
func onApplicationDidBecomeActive() {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Kind: DidBecomeActive,
	}
	defaultApp.lifecycleEvents <- event
}

//export onApplicationWillTerminate
func onApplicationWillTerminate() {
	defaultApp.initWG.Wait()

	event := LifecycleEvent{
		Kind: WillTerminate,
	}
	defaultApp.lifecycleEvents <- event
}

type Orientation int32

const (
	OrientationUnknown Orientation = iota
	// OrientationPortrait when device oriented vertically, home button on the bottom.
	OrientationPortrait
	// OrientationPortraitUpsideDown when device oriented vertically, home button on the top.
	OrientationPortraitUpsideDown
	// OrientationLandscapeLeft when device oriented horizontally, home button on the right.
	OrientationLandscapeLeft
	// OrientationLandscapeRight when device oriented horizontally, home button on the left.
	OrientationLandscapeRight
	// OrientationFaceUp when device oriented flat, face up.
	OrientationFaceUp
	// OrientationFaceDown when device oriented flat, face down.
	OrientationFaceDown
)

type ConfigurationEvent struct {
	NativeWidth  int32
	NativeHeight int32
	NativeScale  float32
	Orientation  Orientation
}

//export onConfigurationChanged
func onConfigurationChanged(w, h int32, scale float32, orientation int32) {
	defaultApp.initWG.Wait()

	out := defaultApp.getConfigurationEventsOut()
	if out == nil {
		return
	}
	event := ConfigurationEvent{
		NativeWidth:  w,
		NativeHeight: h,
		NativeScale:  scale,
		Orientation:  Orientation(orientation),
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

type EventType int32

const (
	EventTypeTouches EventType = iota
	EventTypeMotion
	EventTypeRemoteControl
	EventTypePresses
)

type EventSubtype int32

const (
	// available in iPhone OS 3.0
	EventSubtypeNone EventSubtype = 0

	// for UIEventTypeMotion, available in iPhone OS 3.0
	EventSubtypeMotionShake EventSubtype = 1

	// for UIEventTypeRemoteControl, available in iOS 4.0
	EventSubtypeRemoteControlPlay                 EventSubtype = 100
	EventSubtypeRemoteControlPause                EventSubtype = 101
	EventSubtypeRemoteControlStop                 EventSubtype = 102
	EventSubtypeRemoteControlTogglePlayPause      EventSubtype = 103
	EventSubtypeRemoteControlNextTrack            EventSubtype = 104
	EventSubtypeRemoteControlPreviousTrack        EventSubtype = 105
	EventSubtypeRemoteControlBeginSeekingBackward EventSubtype = 106
	EventSubtypeRemoteControlEndSeekingBackward   EventSubtype = 107
	EventSubtypeRemoteControlBeginSeekingForward  EventSubtype = 108
	EventSubtypeRemoteControlEndSeekingForward    EventSubtype = 109
)

type TouchesState int32

const (
	// TouchesBegan is sent when one or more fingers touch down in a view or window.
	TouchesBegan TouchesState = iota
	// TouchesMoved is sent when one or more fingers associated with an event move within a view or window.
	TouchesMoved
	// TouchesEnded is sent when one or more fingers are raised from a view or window.
	TouchesEnded
	// TouchesCancelled is sent when a system event (such as a low-memory warning) cancels a touch event.
	TouchesCancelled
)

func (state TouchesState) String() string {
	switch state {
	case TouchesBegan:
		return "began"
	case TouchesMoved:
		return "moved"
	case TouchesEnded:
		return "ended"
	case TouchesCancelled:
		return "cancelled"
	default:
		return ""
	}
}

type MotionState int32

const (
	// MotionBegan tells that a motion event has begun.
	MotionBegan MotionState = iota
	// MotionEnded tells that a motion event has ended.
	MotionEnded
	// MotionCancelled tells that a motion event has been cancelled.
	MotionCancelled
)

type TouchSequence int32

type TouchEvent struct {
	State    TouchesState
	Sequence TouchSequence
	X, Y     float32
}

//export onTouchEvent
func onTouchEvent(tp uintptr, state int32, x, y float32) {
	log.Println("[DEBUG touch]", tp, state, x, y)
	defaultApp.initWG.Wait()
	out := defaultApp.getTouchEventsOut()
	if out == nil {
		return
	}

	seq := -1
	for i, val := range touchIDs {
		if val == tp {
			seq = i
			break
		}
	}
	if seq == -1 {
		for i, val := range touchIDs {
			if val == 0 {
				touchIDs[i] = tp
				seq = i
				break
			}
		}
		if seq == -1 {
			panic("maximum touch sequence length exceeded")
		}
	}

	s := TouchesState(state)
	if s == TouchesEnded || s == TouchesCancelled {
		touchIDs[seq] = 0
	}
	event := TouchEvent{
		X:        x,
		Y:        y,
		State:    s,
		Sequence: TouchSequence(seq),
	}
	select {
	case out <- event:
		// dispatched
	case <-time.After(defaultApp.maxDispatchTime):
		// timed out
	}
}

// touchIDs is the current active touches. The position in the array
// is the ID, the value is the UITouch* pointer value.
var touchIDs [11]uintptr
