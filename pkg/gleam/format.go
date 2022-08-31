package gleam

import (
	"bytes"
	"fmt"
	"os/exec"

	pgs "github.com/lyft/protoc-gen-star"
)

func NewGleamFormatter() pgs.PostProcessor { return GleamFormatter{} }

type GleamFormatter struct{}

func (gf GleamFormatter) Match(a pgs.Artifact) bool {
	switch a.(type) {
	case pgs.GeneratorFile, pgs.GeneratorTemplateFile,
		pgs.CustomFile, pgs.CustomTemplateFile:
		return true
	default:
		return false
	}
}

func (gf GleamFormatter) Process(in []byte) ([]byte, error) {
	cmd := exec.Command("/home/sdancer/bin/gleam", "format", "--stdin")
	cmd.Stdin = bytes.NewReader(in)
	stdout, err := cmd.Output()
	if err != nil {
          return in, nil
		return nil, fmt.Errorf("Error formating generated gleam code: %s / %s", string(stdout), err)
	}

	return stdout, nil
}
