package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/thaddock/modrinth-query/api"
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

func download(file api.GetVersionFileInfo, dest string) error {
	destFilename := path.Join(dest, file.Filename)
	fmt.Printf("\tDownloading from: %s to %s\n", file.Url, destFilename)
	cmd := exec.Command("/usr/bin/curl", "-v", "-o", destFilename, file.Url)
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	if _, err := os.Stat(destFilename); err != nil {
		return errors.New(fmt.Sprintf("Unable to download file to: %s", destFilename))
	}
	return nil
}

func latest_versions(c *api.Client, project_id string, loaders []string, game_versions []string) (map[string]api.GetVersionResult, error) {
	vr, err := api.GetVersion(c, project_id, loaders, game_versions, nil)
	if err != nil {
		return nil, err
	}
	wanted_game_versions := make(map[string]bool)
	for _, v := range game_versions {
		wanted_game_versions[v] = true
	}
	res := make(map[string]api.GetVersionResult)
	for _, v := range vr {
		for _, gv := range v.GameVersions {
			if !wanted_game_versions[gv] {
				continue
			}
			cur, ok := res[gv]
			if !ok {
				res[gv] = v
			} else if v.DatePublished > cur.DatePublished {
				res[gv] = v
			}
		}
	}
	return res, nil
}

func DoSync(c *api.Client, dir string, loaders []string, game_versions []string, update bool) error {
	if update && (len(loaders) != 1 || len(game_versions) != 1) {
		return errors.New("Update requires a single loader and game_version")
	}
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
			lv, err := latest_versions(c, vr.ProjectId, loaders, game_versions)
			if err != nil {
				fmt.Printf("\tVERSION ERROR: %s\n", err)
			} else {
				for ver, v := range lv {
					fmt.Printf("\t%s: %s\n", ver, v.Name)
				}
			}
			if update {
				latest, ok := lv[game_versions[0]]
				if !ok {
					fmt.Printf("\t(Not updated, game version not found)\n")
					continue
				}
				if vr.DatePublished == latest.DatePublished {
					fmt.Printf("\tModule up to date\n")
					continue
				}
				fmt.Printf("\tCurrent published: %s / Latest: %s\n", vr.DatePublished, latest.DatePublished)
				if len(latest.Files) != 1 {
					return errors.New("Only support one files for now")
				}
				if err := download(latest.Files[0], dir); err != nil {
					return err
				}
				if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
