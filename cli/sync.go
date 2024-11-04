package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lazerdye/modrinth_query/api"
)

func sha1sum(file string) (string, error) {
	cmd := exec.Command("sha1sum", file)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	s := strings.Split(string(out), " ")
	return s[0], nil
}

func DoSync(c *api.Client, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.Type().IsDir() {
			sha1sum, err := sha1sum(filepath.Join(dir, e.Name()))
			if err != nil {
				fmt.Printf("%s: err: %s\n", e.Name(), err)
				continue
			}
			vr, err := api.GetVersionFile(c, sha1sum, "sha1", false)
			if err != nil {
				fmt.Printf("%s: err: %s\n", e.Name(), err)
				continue
			}
			if vr == nil {
				fmt.Printf("%s: NOT FOUND\n", e.Name())
				continue
			}
			fmt.Printf("%s: %s\n", e.Name(), vr.Name)
		}
	}
	return nil
}
