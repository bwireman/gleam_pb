package gleam

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type gpbWrapper struct {
	pathToBinary string
}

func newGPBWrapper(pathToBinary string) (*gpbWrapper, error) {
	if !exists(pathToBinary) {
		return nil, fmt.Errorf("protoc-erl could not be found at %s", pathToBinary)
	}

	return &gpbWrapper{
		pathToBinary: pathToBinary,
	}, nil
}

func (g *gpbWrapper) generate(targets []string, outputPath string) (err error) {
	args := []string{"-pkgs", "-modname", "gleam_gpb", "-I", ".", "-o", outputPath}
	args = append(args, targets...)

	cmd := exec.Command(g.pathToBinary, args...)

	o, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running protoc-erl: %s", string(o))
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Open(path)
	return !errors.Is(err, os.ErrNotExist)
}
