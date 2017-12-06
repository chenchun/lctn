package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
)

var (
	RootDir = flag.String("root", "", "the root directory of container")
	Init    = flag.Bool("init", false, "is it the init process of container. (This flag is used by lctn internally)")
)

func InitFlags() {
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [Flags] Command [Argument ...]:\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
Command:
  The init command of container

Argument:
  The Arguments of container

Examples:
  lctn -logtostderr -root ./rootfs /bin/sh`)
	}
	flag.Parse()
	if *RootDir == "" {
		*RootDir = "."
	}
	tailArgs := flag.Args()
	if len(tailArgs) < 1 {
		glog.Fatal("command required")
	}
	glog.V(2).Infof("root %s, command %v", *RootDir, tailArgs)
}
