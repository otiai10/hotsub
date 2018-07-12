package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

var commands = []cli.Command{
	quickguide,
}

var quickguide = cli.Command{
	Name:    "quickguide",
	Aliases: []string{"qg"},
	Action: func(ctx *cli.Context) error {

		speak("Hello! This is a quick guide to know how to use hotsub.")
		speak("First, let's check if you have enough tools on your machine.")

		ng := 0

		if dkm, err := exec.LookPath("docker-machine"); err == nil {
			speak("✔\tThe command `docker-machine` is found in %v", dkm)
		} else {
			ng++
			speak("NG!\tdocker-machine: %v", err)
			speak("\tYou need to install `docker-machine`.")
			speak("\tIt's included in `Docker toolbox`, please go to https://docs.docker.com/toolbox/overview/ to install toolbox!")
		}

		speak("\nThen, let's check credentials so that you can access AWS or any cloud services.")

		if fpath, err := checkAWSCredentials(); err == nil {
			speak("✔\tThe AWS credential file is found at %v", fpath)
		} else {
			ng++
			speak("NG!\tAWS Credentials: %v", err)
		}

		if ng == 0 {
			speak("\nCongrats! It seems you are ready to use `hotsub`.")
			speak("For the next step let's try following command.")
			fmt.Println(ocrExample)
			speak("If you have any question for using hotsub, just issue `hotsub help`, or")
			speak("create new issue on https://github.com/otiai10/hotsub/issues. Thank you!")
		}

		return nil
	},
}

func checkAWSCredentials() (string, error) {
	hd, err := homedir()
	if err != nil {
		return "", err
	}
	fpath := filepath.Join(hd, ".aws", "credentials")
	stat, err := os.Stat(fpath)
	if err != nil {
		return "", fmt.Errorf("credential file can't be found at %v: %v", fpath, err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("credential file path seems to be a directory: %v", fpath)
	}
	return fpath, nil
}

func homedir() (string, error) {
	if hd := os.Getenv("HOME"); hd != "" {
		return hd, nil
	}
	myself, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("can't detect current user: %v", err)
	}
	return myself.HomeDir, nil
}

func speak(format string, v ...interface{}) {
	dummysleep := time.Duration(rand.Intn(1000))
	time.Sleep(dummysleep * time.Millisecond)
	format += "\n"
	fmt.Printf(format, v...)
	rand.Seed(time.Now().UnixNano())
}

var ocrExample = `

    /*
     * If you want to try hotsub on AWS, you need 2 things beforehand.
		 * 1) S3 Bucket, in which your input files/directories are located
		 * 2) TODO: more friendly quickstart
     */

    // Prepare your input files on your s3 bucket
    $ aws s3 cp --recursive ./examples/wordcount/speech s3://{YOUR_S3_BUCKET}/speech

    // Edit parameter file to use your s3 bucket
    $ sed -e "s/_placeholder_/{YOUR_S3_BUCKET}/g" ./examples/wordcount/template.csv > ./examples/wordcount/wordcount.csv

    // Execute the tasks with specific docker image
    $ hotsub --tasks ./examples/wordcount/wordcount.csv --script ./examples/wordcount/main.sh --verbose

    /*
     * Then you can see s3://{YOUR_S3_BUCKET}/speech/out is created
     * and the word-count result files.
     */
`
