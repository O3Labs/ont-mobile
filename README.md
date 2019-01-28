# ont-mobile
Useful functions for the Ontology blockchain written in Go.

Built for using with [neo-utils](https://github.com/O3Labs/neo-utils).

### Installation
```
go get github.com/O3Labs/ont-mobile/ontmobile
```

## Compile this library to native mobile frameworks.

### Install gomobile
`go get golang.org/x/mobile/cmd/gomobile`  
`gomobile init`  

### Install Dependencies
```
. dep.sh
```

### Compile to both iOS and Android frameworks
```
. build.sh
```

### Compile to iOS framework
XCode is required.  
`gomobile bind -target=ios -o=output/ios/ontmobile.framework github.com/o3labs/ont-mobile/ontmobile`

### Compile to Android framework
Android NDK is required. https://developer.android.com/ndk/guides/index.html  
```
gomobile init -ndk ~/Library/Android/sdk/ndk-bundle/

ANDROID_HOME=/Users/$USER/Library/Android/sdk gomobile bind -target=android -o=output/android/ontmobile.aar github.com/o3labs/ont-mobile/ontmobile
```
