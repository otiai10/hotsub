package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	headExpression = regexp.MustCompile("^(?P<key>.+) +(?P<bind>.+)$")
)

// Task ...
type Task struct {
	Env             map[string]string
	Inputs          map[string]string
	InputRecursive  map[string]string
	Outputs         map[string]string
	OutputRecursive map[string]string
	Index           int
}

// NewTask ...
func NewTask(i ...int) *Task {
	if len(i) == 0 {
		i = append(i, 0)
	}
	return &Task{
		Index:           i[0],
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
	ext := filepath.Ext(taskfile.Name())
	switch ext {
	case ".csv":
		r := csv.NewReader(taskfile)
		return parseTasksFromRowReader(r)
	case ".tsv":
		r := csv.NewReader(taskfile)
		r.Comma = '\t'
		return parseTasksFromRowReader(r)
	default:
		return nil, fmt.Errorf("unexpected extension for task file: %v", ext)
	}
}

func parseTasksFromRowReader(r *csv.Reader) ([]*Task, error) {
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
		task := NewTask(i)
		for colindex, value := range row {
			if err := colums[colindex].Bind(task, value); err != nil {
				return nil, err
			}
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
