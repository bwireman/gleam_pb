package main

import (
	"github.com/bwireman/gleam_pb/pkg/gleam"
	pgs "github.com/lyft/protoc-gen-star"
)

func main() {
	pgs.Init(pgs.DebugEnv("DEBUG")).RegisterModule(
		gleam.Gleam(),
	).Render()
}
