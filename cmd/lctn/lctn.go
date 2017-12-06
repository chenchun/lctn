package main

import (
	"github.com/chenchun/lctn"
	"github.com/chenchun/lctn/flags"
	"github.com/golang/glog"
)

func main() {
	flags.InitFlags()
	defer glog.Flush()
	if !*flags.Init {
		lctn.Parent()
	} else {
		lctn.Child()
	}
}
