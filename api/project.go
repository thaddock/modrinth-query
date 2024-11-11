package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Project struct {
	ProjectId         string   `json:"project_id"`
	ProjectType       string   `json:"project_type"`
	Slug              string   `json:"slug"`
	Description       string   `json:"description"`
	Categories        []string `json:"categories"`
	DisplayCategories []string `json:"display_categories"`
	Versions          []string `json:"versions"`
	Downloads         int64    `json:"downloads"`
	Follows           int64    `json:"follows"`
	IconUrl           string   `json:"icon_url"`
	DateCreated       string   `json:"date_created"`
	DateModified      string   `json:"date_modified"`
	LatestVersion     string   `json:"latest_version"`
	License           string   `json:"license"`
	ClientSide        string   `json:"client_side"`
	ServerSide        string   `json:"server_side"`
	Gallery           []string `json:"gallery"`
	FeaturedGallery   string   `json:"featured_gallery"`
	Color             int64    `json:"color"`
}

type SearchResult struct {
	Hits      []Project `json:"hits"`
	Offset    int32     `json:"offset"`
	Limit     int32     `json:"limit"`
	TotalHits int32     `json:"total_hits"`
}

func ProjectSearch(c *Client, query string, facets string, index string, offset int32, limit int32) (*SearchResult, error) {
	params := url.Values{}
	params.Add("query", query)
	if facets != "" {
		params.Add("facets", facets)
	}
	if index != "" {
		params.Add("index", index)
	}
	params.Add("offset", strconv.Itoa(int(offset)))
	params.Add("limit", strconv.Itoa(int(limit)))
	url := fmt.Sprintf("%s/search?%s", c.BaseUrl, params.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sr SearchResult
	err = json.Unmarshal(body, &sr)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}
