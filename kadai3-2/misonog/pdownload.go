package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"runtime"
)

// Pdownload structs
type Pdownload struct {
	Utils
	URL       string
	TargetDir string
	Procs     int
	useragent string
	referer   string
}

func New() *Pdownload {
	return &Pdownload{
		Utils: &Data{},
		Procs: runtime.NumCPU(), // default
	}
}

func (pdownload *Pdownload) Run() error {
	if err := pdownload.Ready(); err != nil {
		return err
	}

	if err := pdownload.Check(); err != nil {
		return err
	}

	if err := pdownload.Download(); err != nil {
		return err
	}

	if err := pdownload.Utils.MergeFiles(pdownload.Procs); err != nil {
		return err
	}

	return nil
}

func (pdownload *Pdownload) Ready() error {
	var targetDir string

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	flag.StringVar(&targetDir, "d", pwd, "path to the directory to save the downloaded file, filename will be taken from url")
	flag.Parse()

	if err := pdownload.parseURL(flag.Args()); err != nil {
		return err
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("target directory is not exist: %v", err)
	}
	pdownload.TargetDir = targetDir

	return nil
}

func (pdownload *Pdownload) parseURL(args []string) error {
	if len(args) > 1 {
		return errors.New("URL must be a single")
	}
	if len(args) < 1 {
		return errors.New("urls not found in the arguments passed")
	}

	for _, arg := range args {
		_, err := url.ParseRequestURI(arg)
		if err != nil {
			return err
		}
		pdownload.URL = arg
	}

	return nil
}
