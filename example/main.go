package main

import (
	"log"

	"github.com/xlab/ios-go/app"
)

func main() {
	log.Println("GoApp has started ^_^")
	configEvents := make(chan app.ConfigurationEvent, 1)
	touchEvents := make(chan app.TouchEvent, 10)

	app.Main(func(a app.AppDelegate) {
		a.HandleConfigurationEvents(configEvents)
		a.HandleTouchEvents(touchEvents)
		a.InitDone()
		for {
			select {
			case event := <-a.LifecycleEvents():
				switch event.Kind {
				case app.ViewDidLoad:
					log.Println(event.Kind, "handled")
				default:
					log.Println(event.Kind, "event ignored")
				}
			case cfg := <-configEvents:
				log.Printf("rotated device: %+v\n", cfg)
			case tc := <-touchEvents:
				log.Printf("touch[%d]: (%.1f,%.1f) -> %s\n", tc.Sequence, tc.X, tc.Y, tc.State)
			case <-a.VSync():
				// no-op
			}
		}
	})
}
