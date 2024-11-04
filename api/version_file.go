package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GetVersionFileInfo struct {
	Hashes   map[string]string `json:"hashes"`
	Url      string            `json:"url"`
	Filename string            `json:"filename"`
	Primary  bool              `json:"primary"`
	Size     int               `json:"size"`
	FileType string            `json:"file_type"`
}

type GetVersionDependency struct {
	VersionId      string `json:"version_id"`
	ProjectId      string `json:"project_id"`
	FileName       string `json:"file_name"`
	DependencyType string `json:"dependency_type"`
}

type GetVersionFileResult struct {
	GameVersions    []string               `json:"game_versions"`
	Loaders         []string               `json:"loaders"`
	Id              string                 `json:"id"`
	ProjectId       string                 `json:"project_id"`
	AuthorId        string                 `json:"author_id"`
	Featured        bool                   `json:"featured"`
	Name            string                 `json:"name"`
	VersionNumber   string                 `json:"version_number"`
	Changelog       string                 `json:"changelog"`
	ChangelogUrl    string                 `json:"changelog_url"`
	DatePublished   string                 `json:"date_published"`
	Downloads       int                    `json:"downloads"`
	VersionType     string                 `json:"version_type"`
	Status          string                 `json:"status"`
	RequestedStatus string                 `json:"requested_status"`
	Files           []GetVersionFileInfo   `json:"files"`
	Dependencies    []GetVersionDependency `json:"dependencies"`
}

func GetVersionFile(c *Client, hash string, algorithm string, multiple bool) (*GetVersionFileResult, error) {
	params := url.Values{}
	params.Add("algorithm", algorithm)
	if multiple {
		params.Add("multiple", "true")
	}
	url := fmt.Sprintf("%s/version_file/%s?%s", QueryUrl, hash, params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var vr GetVersionFileResult
	err = json.Unmarshal(body, &vr)
	if err != nil {
		return nil, err
	}
	return &vr, nil
}
