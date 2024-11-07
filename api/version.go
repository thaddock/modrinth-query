package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GetVersionResult struct {
	Name            string                 `json:"name"`
	VersionNumber   string                 `json:"version_number"`
	Changelog       string                 `json:"changelog"`
	Dependencies    []GetVersionDependency `json:"dependencies"`
	GameVersions    []string               `json:"game_versions"`
	VersionType     string                 `json:"version_type"`
	Loaders         []string               `json:"loaders"`
	Featured        bool                   `json:"featured"`
	Status          string                 `json:"status"`
	RequestedStatus string                 `json:"requested_status"`
	Id              string                 `json:"id"`
	ProjectId       string                 `json:"project_id"`
	AuthorId        string                 `json:"author_id"`
	DatePublished   string                 `json:"date_published"`
	Downloads       int64                  `json:"downloads"`
	ChangelogUri    string                 `json:"changelog_uri"`
	Files           []GetVersionFileInfo   `json:"files"`
}

func GetVersion(c *Client, project_id string, loaders []string, game_versions []string, featured *bool) ([]GetVersionResult, error) {
	params := url.Values{}
	if len(loaders) > 0 {
		params.Add("loaders", "[\""+strings.Join(loaders, "\",\"")+"\"]")
	}
	if len(game_versions) > 0 {
		params.Add("game_versions", "[\""+strings.Join(game_versions, "\",\"")+"\"]")
	}
	if featured != nil && *featured {
		params.Add("featured", "true")
	} else if featured != nil && !*featured {
		params.Add("featured", "false")
	}
	url := fmt.Sprintf("%s/project/%s/version?%s", QueryUrl, project_id, params.Encode())
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
	var vr []GetVersionResult
	err = json.Unmarshal(body, &vr)
	if err != nil {
		return nil, err
	}
	return vr, nil
}
