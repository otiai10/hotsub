package application

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/otiai10/hotsub/params"
)

// Init inits cli environment on which hotsub runs.
func Init(ctx params.Context) error {

	if err := dkmsn(); err != nil {
		return err
	}

	if err := awsconfig(); err != nil {
		return err
	}

	fmt.Println("Prerequisites check [DONE]")
	return nil
}

func dkmsn() error {
	name := "docker-machine"
	if _, err := exec.LookPath(name); err != nil {
		if execerr, ok := err.(*exec.Error); ok && execerr.Err == exec.ErrNotFound {
			fmt.Println(`"docker-machine" is not found.`)
			fmt.Println("\tGo to https://docs.docker.com/machine/install-machine/")
		} else {
			return err
		}
	}
	return nil
}

func awsconfig() error {
	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".aws", "credentials")); err != nil {
		if os.IsNotExist(err) {
			fmt.Println(`".aws/credentials" file is not located under $HOME.`)
			fmt.Println("\tGo to https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html")
		} else {
			return err
		}
	}
	return nil
}
