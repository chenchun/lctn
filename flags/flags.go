package flags

import (
	"flag"

	"github.com/golang/glog"
)

var (
	RootDir = flag.String("root", "", "the root directory of container")
)

func InitFlags() {
	flag.Parse()
	defer glog.Flush()
	if *RootDir == "" {
		*RootDir = "."
	}
	tailArgs := flag.Args()
	if len(tailArgs) < 1 {
		glog.Fatal("command required")
	}
	glog.V(1).Infof("root %s, command %v", *RootDir, tailArgs)
}
