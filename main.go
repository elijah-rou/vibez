package main

import (
	"fmt"
	"os"

	"github.com/simone-vibes/vibez/cmd"
	"github.com/simone-vibes/vibez/internal/crash"
	"github.com/simone-vibes/vibez/internal/version"
)

func main() {
	defer crash.Recover("main")
	if err := crash.Install(version.Version); err != nil {
		fmt.Fprintf(os.Stderr, "vibez: crash logging unavailable: %v\n", err)
	}
	cmd.Execute()
}
