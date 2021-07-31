package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// execute cmd
func (i Item) Render(d *InitData) error {
	// make proxy dir
	if i.Proxy != "" && !Exists(filepath.Dir(i.Proxy)) {
		err := os.MkdirAll(filepath.Dir(i.Proxy), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	// make output Dir
	if !Exists(filepath.Dir(i.Output)) {
		err := os.MkdirAll(filepath.Dir(i.Output), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("command : ", strings.Join(i.FirstCmd, " "))
	cmd := exec.Command(i.FirstCmd[0], i.FirstCmd[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "OCIO="+d.Ocio)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	if i.SecondCmd != nil {
		log.Println("command : ", strings.Join(i.SecondCmd, " "))
		cmd := exec.Command(i.SecondCmd[0], i.SecondCmd[1:]...)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "OCIO="+d.Ocio)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
