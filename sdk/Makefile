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