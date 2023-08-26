// Trugw server c-shared library
//
// Build shared library for linux:
//
//	go build -buildmode c-shared -o libtrugw.so .
//
// Build dll for Windows:
//
//	CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o libtru.dll .
//
// Build 32bit dll for Windows:
//
//	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -o libtru.dll .
//
// Build dll for Windows under Windows:
/*
	set GOOS=windows
	set GOARCH=amd64
	set CGO_ENABLED=1
	set CC=c:\mingw64\bin\x86_64-w64-mingw32-gcc.exe
	set CXX=c:\mingw64\bin\x86_64-w64-mingw32-g++.exe
	set PKG_CONFIG_PATH=c:\mingw64\lib\pkgconfig
	set MSYS2_ARCH=x86_64

	go build -buildmode=c-shared -o libtru.dll .
*/
// or use bat file:
//
//  build.bat
//
package main

// #include <stdlib.h>
//
//
import "C"
import (
	"fmt"
	"os"
)

func main() {}

var sockAddr = os.TempDir() + "/trugw.sock"

func init() {
	server()
}

// Run the trugw server
//
//export server
func server() C.int {
	fmt.Printf("Tru unix socket gateway server, sock path: %s\n", sockAddr)
	return 0
}
