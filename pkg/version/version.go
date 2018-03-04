package version

import (
	"flag"
	"fmt"
	"os"
)

var (
	Build   string
	Version string

	v = flag.Bool("v", false, "Gets version of binary")
)

// CheckVersionFlag checks to see if v was passed and prints out info
func CheckVersionFlag() {
	if *v {
		fmt.Println("Version:", Version)
		fmt.Println("Build:", Build)
		os.Exit(0)
	}
}
