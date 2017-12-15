package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	headExpression = regexp.MustCompile("^(?P<key>.+) +(?P<bind>.+)$")
)

// Entry ...
type Entry struct {
	URL      string
	Filepath string
}

// Task ...
type Task struct {
	Index           int
	Prefix          string
	Env             map[string]string
	Inputs          map[string]string
	InputRecursive  map[string]string
	Outputs         map[string]string
	OutputRecursive map[string]string

	// ContainerEnv is a set of real and flattened env variables of
	// Env, Inputs, InputRecursive, Outputs, and OutputRecursive.
	// These values are supposed to be created by Handler.Prepare lifecycle.
	ContainerEnv []string
}

// NewTask ...
func NewTask(i int, prefix string) *Task {
	return &Task{
		Index:           i,
		Prefix:          prefix,
		Env:             map[string]string{},
		Inputs:          map[string]string{},
		InputRecursive:  map[string]string{},
		Outputs:         map[string]string{},
		OutputRecursive: map[string]string{},
	}
}

// Column ...
type Column struct {
	Type string
	Name string
}

// Bind ...
func (c Column) Bind(task *Task, value string) error {
	switch c.Type {
	case "--env":
		task.Env[c.Name] = value
	case "--input":
		task.Inputs[c.Name] = value
	case "--input-recursive":
		task.InputRecursive[c.Name] = value
	case "--output":
		task.Outputs[c.Name] = value
	case "--output-recursive":
		task.OutputRecursive[c.Name] = value
	}
	return nil
}

func parseTasksFromFile(taskfile *os.File) ([]*Task, error) {
	name := strings.TrimRight(filepath.Base(taskfile.Name()), filepath.Ext(taskfile.Name()))
	ext := filepath.Ext(taskfile.Name())
	switch ext {
	case ".csv":
		r := csv.NewReader(taskfile)
		return parseTasksFromRowReader(r, name)
	case ".tsv":
		r := csv.NewReader(taskfile)
		r.Comma = '\t'
		return parseTasksFromRowReader(r, name)
	default:
		return nil, fmt.Errorf("unexpected extension for task file: %v", ext)
	}
}

func parseTasksFromRowReader(r *csv.Reader, prefix string) ([]*Task, error) {
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	header, rows := rows[0], rows[1:]

	tasks := []*Task{}
	colums := []Column{}
	for _, th := range header {
		matched := headExpression.FindStringSubmatch(th)
		if len(matched) < 3 {
			return nil, fmt.Errorf("unexpected format for task file columns header: %v", th)
		}
		colums = append(colums, Column{Type: matched[1], Name: matched[2]})
	}
	for i, row := range rows {
		if len(row) < len(header) {
			return tasks, fmt.Errorf("csv/tsv record doesn't have enough columns specified with the first row: %v", i)
		}
		task := NewTask(i, prefix)
		for colindex, value := range row {
			if err := colums[colindex].Bind(task, value); err != nil {
				return nil, err
			}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
