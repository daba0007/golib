package golib

import (
	"fmt"
)

type golib struct {
	ver   string
	gover string
}

var (
	// Golib ...
	Golib = &golib{
		ver:   "1.0",
		gover: "1.13+",
	}
)

func (s *golib) Description() string {
	return fmt.Sprintf("golib verion %s, build on go %s", s.ver, s.gover)
}
