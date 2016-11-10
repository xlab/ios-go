package app

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework QuartzCore

#include <UIKit/UIDevice.h>

void runApp(void);
*/
import "C"
import "runtime"

func init() {
	// Lock the goroutine responsible for initialization to an OS thread.
	// This means the goroutine running main (and calling the run function
	// below) is locked to the OS thread that started the program. This is
	// necessary for the correct delivery of UIKit events to the process.
	//
	// A discussion on this topic:
	// https://groups.google.com/forum/#!msg/golang-nuts/IiWZ2hUuLDA/SNKYYZBelsYJ
	runtime.LockOSThread()

	defaultApp.initWG.Add(1)
}

func Main(f func(AppDelegate)) {
	defer runtime.UnlockOSThread()

	go f(defaultApp) // run in a separate thread
	C.runApp()       // remains bound to the OS thread
	panic("runApp unexpected exit")
}
