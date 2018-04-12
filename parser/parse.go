package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otiai10/awsub/core"
)

var (
	headExpression   = regexp.MustCompile("^(?P<key>.+) +(?P<bind>.+)$")
	sharedExpression = regexp.MustCompile("^(?P<name>[0-9A-Z_]+)=(?P<url>.+)$")
)

// ParseFile ...
func ParseFile(fpath string) (jobs []*core.Job, err error) {
	abspath, err := filepath.Abs(fpath)
	if err != nil {
		return []*core.Job{}, err
	}
	ext := filepath.Ext(fpath)
	name := strings.TrimRight(filepath.Base(fpath), ext)
	f, err := os.Open(abspath)
	if err != nil {
		return []*core.Job{}, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	switch ext {
	case ".csv":
		return ParseRowReader(r, name)
	case ".tsv":
		r.Comma = '\t'
		r.LazyQuotes = true
		return ParseRowReader(r, name)
	default:
		return nil, fmt.Errorf("unexpected extension for task file: %v", ext)
	}
}

// ParseRowReader ...
func ParseRowReader(r *csv.Reader, prefix string) (jobs []*core.Job, err error) {

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// If this file is empty, return without errors
	if len(rows) == 0 {
		return []*core.Job{}, nil
	}

	hrow, rows := rows[0], rows[1:]

	header := []Column{}
	for _, th := range hrow {
		matched := headExpression.FindStringSubmatch(th)
		if len(matched) < 3 {
			return nil, fmt.Errorf("unexpected format for task file columns header: %v", th)
		}
		header = append(header, Column{
			Type: strings.Trim(matched[1], " "),
			Name: strings.Trim(matched[2], " "),
		})
	}
	for i, row := range rows {
		if len(row) < len(header) {
			return jobs, fmt.Errorf("csv/tsv record doesn't have enough columns specified with the first row: %v", i)
		}
		job := core.NewJob(i, prefix)
		for i, value := range row {
			if err := header[i].Bind(job, value); err != nil {
				return nil, err
			}
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Column ...
type Column struct {
	Type string
	Name string
}

// Bind ...
func (c Column) Bind(job *core.Job, value string) error {
	value = strings.Trim(value, " ")
	switch c.Type {
	case "--env":
		job.Parameters.Envs = append(job.Parameters.Envs, core.Env{Name: c.Name, Value: value})
	case "--input":
		job.Parameters.Inputs = append(
			job.Parameters.Inputs,
			&core.Input{Resource: core.Resource{Name: c.Name, URL: value}},
		)
	case "--input-recursive":
		job.Parameters.Inputs = append(
			job.Parameters.Inputs,
			&core.Input{Resource: core.Resource{Name: c.Name, URL: value, Recursive: true}},
		)
	case "--output":
		job.Parameters.Outputs = append(
			job.Parameters.Outputs,
			&core.Output{Resource: core.Resource{Name: c.Name, URL: value}},
		)
	case "--output-recursive":
		job.Parameters.Outputs = append(
			job.Parameters.Outputs,
			&core.Output{Resource: core.Resource{Name: c.Name, URL: value, Recursive: true}},
		)
	}
	return nil
}

// ParseSharedData ...
func ParseSharedData(kvpairs []string) (inputs core.Inputs, err error) {
	if len(kvpairs) == 0 {
		return
	}
	for _, kv := range kvpairs {
		kvl := sharedExpression.FindStringSubmatch(kv)
		if len(kvl) < 3 {
			err = fmt.Errorf("Invalid format for shared data: %s", kv)
			return
		}
		inputs = append(inputs, &core.Input{Resource: core.Resource{
			Name:      kvl[1],
			URL:       kvl[2],
			Recursive: true, // TEMP
		}})
	}
	return
}
