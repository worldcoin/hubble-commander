# go/Makefile

ios-arm64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=arm64 \
	SDK=iphoneos \
	CC=$(PWD)/clangwrap.sh \
	CGO_CFLAGS="-fembed-bitcode" \
	go build -ldflags "-s -w" -buildmode=c-archive -tags ios -o $(IOS_OUT)/arm64.a main.go

ios-x86_64:
	CGO_ENABLED=1 \
	GOOS=darwin \
	GOARCH=amd64 \
	SDK=iphonesimulator \
	CC=$(PWD)/clangwrap.sh \
	go build -ldflags "-s -w" -buildmode=c-archive -tags ios -o $(IOS_OUT)/x86_64.a main.go

ios: ios-arm64 ios-x86_64
	lipo $(IOS_OUT)/x86_64.a $(IOS_OUT)/arm64.a -create -output $(IOS_OUT)/sdk.a
	cp $(IOS_OUT)/arm64.h $(IOS_OUT)/sdk.h

# Or configured up by build system. Then copy to app/src/main/jniLibs
ANDROID_OUT=./jniLibs
ANDROID_SDK=$(HOME)/Library/Android/sdk
# Put your NDK version or configure by build system
NDK_VERSION=22.1.7171670
NDK_BIN=$(ANDROID_SDK)/ndk/$(NDK_VERSION)/toolchains/llvm/prebuilt/darwin-x86_64/bin

android-armv7a:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm \
	GOARM=7 \
	CC=$(NDK_BIN)/armv7a-linux-androideabi21-clang \
	go build -ldflags "-s -w" -buildmode=c-shared -o $(ANDROID_OUT)/armeabi-v7a/libsdk.so ./main.go

android-arm64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CC=$(NDK_BIN)/aarch64-linux-android21-clang \
	go build -ldflags "-s -w" -buildmode=c-shared -o $(ANDROID_OUT)/arm64-v8a/libsdk.so ./main.go

android-x86:
	CGO_ENABLED=1 \
	CGO_CFLAGS="-mno-stackrealign" \
	GOOS=android \
	GOARCH=386 \
	CC=$(NDK_BIN)/i686-linux-android21-clang \
	go build -ldflags "-s -w" -buildmode=c-shared -o $(ANDROID_OUT)/x86/libsdk.so ./main.go

android-x86_64:
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=$(NDK_BIN)/x86_64-linux-android21-clang \
	go build -ldflags "-s -w" -buildmode=c-shared -o $(ANDROID_OUT)/x86_64/libsdk.so ./main.go

android: android-armv7a android-arm64 android-x86 android-x86_64
