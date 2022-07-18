package pkg

import (
	"os"
)

func OtherFunc() {
	// This function should not trigger analyzer
	os.Exit(1)
}

func main() {
	os.Exit(1) // want `has os.Exit function`
}
