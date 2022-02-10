package gleam

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type gpbWrapper struct {
	pathToBinary string
	outputPath   string
}

func newGPBWrapper(pathToBinary string, outputPath string) (*gpbWrapper, error) {
	if _, exists := exists(pathToBinary); !exists {
		return nil, fmt.Errorf("protoc-erl could not be found at %s", pathToBinary)
	}

	return &gpbWrapper{
		pathToBinary: pathToBinary,
		outputPath:   outputPath,
	}, nil
}

func (g *gpbWrapper) generate(targets []string) (err error) {
	args := []string{"-pkgs", "-modname", "gleam_gpb", "-I", ".", "-o", g.outputPath}
	args = append(args, targets...)

	cmd := exec.Command(g.pathToBinary, args...)

	o, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running protoc-erl: %s", string(o))
	}

	return nil
}

const headerFileName = "gpb.hrl"

func (g *gpbWrapper) updateImport(correctPath string) (err error) {
	expectedPath := path.Join(g.outputPath, "gleam_gpb.erl")

	if f, exists := exists(expectedPath); !exists {
		return fmt.Errorf("updateImport failed: Could not find %s", expectedPath)
	} else if b, err := ioutil.ReadAll(f); err != nil {
		return err
	} else {
		contents := string(b)

		newHeaderPath := path.Join(correctPath, headerFileName)
		newContents := strings.Replace(contents, formatReplace(headerFileName), formatReplace(newHeaderPath), 1)

		if err := os.WriteFile(expectedPath, []byte(newContents), 0666); err != nil {
			return err
		}
	}

	return nil
}

func formatReplace(path string) string {
	return fmt.Sprintf("-include(\"%s\").", path)
}

func exists(path string) (*os.File, bool) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) || err != nil {
		return nil, false
	}

	return f, true
}
