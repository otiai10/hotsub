package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/otiai10/dkmachine/v0/dkmachine"
	"github.com/urfave/cli"
)

var (
	headExpression = regexp.MustCompile("^(?P<key>.+) +(?P<bind>.+)$")
)

// Task ...
type Task struct {
	Env            map[string]string
	Inputs         map[string]string
	InputRecursive map[string]string
	Index          int
}

// action ...
func action(ctx *cli.Context) error {
	tasksfpath := ctx.String("tasks")
	f, err := os.Open(tasksfpath)
	if err != nil {
		return fmt.Errorf("failed to open tasks file `%s`: %v", tasksfpath, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	switch filepath.Ext(f.Name()) {
	case ".tsv":
		r.Comma = '\t'
	}
	rows, err := r.ReadAll()
	if err != nil {
		return err
	}
	tasks, err := parseTasks(rows[0], rows[1:])
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", tasks)

	return nil
}

func parseTasks(header []string, rows [][]string) ([]*Task, error) {
	tasks := []*Task{}
	// parsersForColumn := []func(task *Task) error{}
	for _, th := range header {
		fmt.Println(th, headExpression.MatchString(th))
		/*
			matched := headExpression.FindStringSubmatch(th)
			fmt.Printf("%+v\n", matched)
		*/
	}
	for i, row := range rows {
		if len(row) < len(header) {
			return tasks, fmt.Errorf("csv/tsv record doesn't have enough columns specified with the first row: %v", i)
		}
	}
	return tasks, nil
}

func handleTask(task string, wg *sync.WaitGroup) {
	co := &dkmachine.CreateOptions{
		Driver:                "amazonec2",
		AmazonEC2Region:       "ap-southeast-2",
		AmazonEC2InstanceType: "m4.large",
	}
	machine, err := dkmachine.Create(co)
	fmt.Printf("%+v\n", machine)
	fmt.Println(err)
	wg.Done()
}
