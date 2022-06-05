package system_info

import (
	"fmt"
	"runtime"
)

// GetSystemInfo function makes sense about the system which runs the app
func GetSystemInfo() {
	fmt.Printf("GOMAXPROC is %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU is %d\n", runtime.NumCPU())
	fmt.Printf("GOARCH is %v\n", runtime.GOARCH)
	fmt.Printf("GOOS is %v\n", runtime.GOOS)
}
