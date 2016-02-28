package goflat

import (
	"io"
	"os/exec"
	"path/filepath"
	"strings"
)

//Flat struct
type Flat struct {
	MainGo     string
	GoTemplate string
	GoInputs   []goInput
	GoPipes    string
}

//GoRun runs go on the dynamically created main.go with a given stdout and stderr pipe
func (f *Flat) GoRun(outWriter io.Writer, errWriter io.Writer) error {
	out := make([]string, len(f.GoInputs)+3)
	out[0], out[1], out[2] = "run", f.MainGo, f.GoPipes
	for k, v := range f.GoInputs {
		out[k+3] = v.Path
	}
	cmd := exec.Command("go", out...)
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	return cmd.Run()
}

type goInput struct{ Path, StructName, VarName string }

func newGoInput(input string) goInput {
	//optionally the structname can be passed via commandline with ":" seperator
	if strings.Contains(input, ":") {
		s := strings.Split(input, ":")
		return goInput{Path: s[0], StructName: s[1]}
	}
	//goflat convention is to build a structname based on filename using strings Title convention
	name := filepath.Base(input)
	name = strings.Title(strings.Split(name, ".")[0])
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, "_", "", -1)
	return goInput{
		Path:       input,
		StructName: name,
	}
}
