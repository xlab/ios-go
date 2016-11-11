WORK ?= $(shell pwd)
PACKAGE ?= github.com/xlab/ios-go/example

CC ?= $(shell xcrun --sdk iphoneos --find clang)
CXX ?= $(shell xcrun --sdk iphoneos --find clang)
ISYSROOT ?= $(shell xcrun --sdk iphoneos --show-sdk-path)

all: build

build:
	mkdir -p $(WORK)/build
	GOOS=darwin GOARCH=arm GOARM=7 CC=$(CC) CXX=$(CXX) CGO_CFLAGS="-isysroot $(ISYSROOT) -miphoneos-version-min=8.0 -arch armv7" CGO_LDFLAGS="-isysroot $(ISYSROOT) -miphoneos-version-min=8.0 -arch armv7" CGO_ENABLED=1 go build -tags="ios" -o=$(WORK)/build/arm $(PACKAGE)
	GOOS=darwin GOARCH=arm64 CC=$(CC) CXX=$(CXX) CGO_CFLAGS="-isysroot $(ISYSROOT) -miphoneos-version-min=8.0 -arch arm64" CGO_LDFLAGS="-isysroot $(ISYSROOT) -miphoneos-version-min=8.0 -arch arm64" CGO_ENABLED=1 go build -tags="ios" -o=$(WORK)/build/arm64 $(PACKAGE)
	xcrun lipo -create $(WORK)/build/arm $(WORK)/build/arm64 -o $(WORK)/main/main
	xcrun xcodebuild -configuration Release -project $(WORK)/main.xcodeproj

install:
	# https://github.com/phonegap/ios-deploy
	ios-deploy -b $(WORK)/build/Release-iphoneos/GoAppExample.app

debug:
	ios-deploy -d -W -b $(WORK)/build/Release-iphoneos/GoAppExample.app

clean:
	rm -rf $(WORK)/build
	rm -f $(WORK)/main/main

.PHONY: build