#import "GoAppViewController.h"
#include "_cgo_export.h"

// #include <MoltenVK/vk_mvk_datatypes.h>
// #include "../Demos.h"

#pragma mark -
#pragma mark GoAppController

@implementation GoAppController {
	CADisplayLink* displayLink;
}

-(void) dealloc {
	[displayLink release];
	[super dealloc];
}

-(void) viewDidLoad {
	[super viewDidLoad];

	self.view.contentScaleFactor = UIScreen.mainScreen.nativeScale;

	// int scale = 1;
	// if ([[UIScreen mainScreen] respondsToSelector:@selector(displayLinkWithTarget:selector:)]) {
	// 	scale = (int)[UIScreen mainScreen].scale; // either 1.0, 2.0, or 3.0.
	// }
	// setScreen(scale);

	CGSize size = [UIScreen mainScreen].nativeBounds.size;
	CGFloat scale = [UIScreen mainScreen].nativeScale;
	UIInterfaceOrientation orientation = [[UIApplication sharedApplication] statusBarOrientation];
	onConfigurationChanged((int)size.width, (int)size.height, scale, orientation);

	onViewDidLoad((GoUintptr)self.view);

	uint32_t fps = 60;
	displayLink = [CADisplayLink displayLinkWithTarget: self selector: @selector(renderLoop)];
	[displayLink setFrameInterval: 60 / fps];
	[displayLink addToRunLoop: NSRunLoop.currentRunLoop forMode: NSDefaultRunLoopMode];
}

-(void) renderLoop {
	onVSync();
}

- (void)viewWillTransitionToSize:(CGSize)size withTransitionCoordinator:(id<UIViewControllerTransitionCoordinator>)coordinator {
	[coordinator animateAlongsideTransition:^(id<UIViewControllerTransitionCoordinatorContext> context) {
		// animate something here
	} completion:^(id<UIViewControllerTransitionCoordinatorContext> context) {
		UIInterfaceOrientation orientation = [[UIApplication sharedApplication] statusBarOrientation];
		CGFloat scale = [UIScreen mainScreen].nativeScale;
		onConfigurationChanged((int)size.width, (int)size.height, scale, orientation);
	}];
}

static void withTouches(int state, NSSet* touches) {
	CGFloat scale = [UIScreen mainScreen].scale;
	for (UITouch* touch in touches) {
		CGPoint p = [touch locationInView:touch.view];
		onTouchEvent((GoUintptr)touch, state, p.x*scale, p.y*scale);
	}
}

- (void)touchesBegan:(NSSet*)touches withEvent:(UIEvent*)event {
	withTouches(0, touches);
}

- (void)touchesMoved:(NSSet*)touches withEvent:(UIEvent*)event {
	withTouches(1, touches);
}

- (void)touchesEnded:(NSSet*)touches withEvent:(UIEvent*)event {
	withTouches(2, touches);
}

- (void)touchesCancelled:(NSSet*)touches withEvent:(UIEvent*)event {
    withTouches(3, touches);
}

@end


#pragma mark -
#pragma mark GoAppView

@implementation GoAppView

/** Returns a Metal-compatible layer. */
+(Class) layerClass { return [CAMetalLayer class]; }

@end
