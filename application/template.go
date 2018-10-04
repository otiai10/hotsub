package application

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/hotsub/params"
)

// Template ...
func Template(ctx params.Context) error {

	name := ctx.String("name")
	dir := ctx.String("dir")
	proj := fmt.Sprintf("%s_%s", name, time.Now().Format("200601021504"))

	dir = filepath.Join(dir, proj)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	tpl := catalog[name]

	fcsv, err := os.Create(filepath.Join(dir, tpl.Tasks.Name))
	if err != nil {
		return err
	}
	defer fcsv.Close()
	wcsv := csv.NewWriter(fcsv)
	if err := wcsv.WriteAll(tpl.Tasks.Table); err != nil {
		return err
	}
	wcsv.Flush()

	fscript, err := os.Create(filepath.Join(dir, tpl.Script.Name))
	if err != nil {
		return err
	}
	defer fscript.Close()
	if _, err := io.WriteString(fscript, tpl.Script.Content); err != nil {
		return err
	}

	freadme, err := os.Create(filepath.Join(dir, tpl.README.Name))
	if err != nil {
		return err
	}
	defer freadme.Close()
	if _, err := io.WriteString(freadme, tpl.README.Content); err != nil {
		return err
	}

	fmt.Printf("Successfully created template project\n    at %v\n", dir)
	fmt.Println("Go to the directory and check README.md")

	return nil
}

type file struct {
	Name    string
	Content string
	Table   [][]string
}

type template struct {
	Script file
	README file
	Tasks  file
}

var catalog = map[string]template{
	"helloworld": template{
		Script: file{
			Name: "hello.sh",
			Content: `#!/bin/bash

echo "Hello! My name is ${NAME}."
echo "I'm a ${SPECIES}."
`,
		},
		README: file{
			Name:    "README.md",
			Content: "hotsub run --script hello.sh --tasks hello.csv --verbose",
		},
		Tasks: file{
			Name: "hello.csv",
			Table: [][]string{
				[]string{"--env NAME", "--env SPECIES"},
				[]string{"Tom", "cat"},
				[]string{"Jerry", "mouse"},
			},
		},
	},
}
