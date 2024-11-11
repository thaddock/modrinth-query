package api

type Client struct {
	BaseUrl string
}

func NewClient() *Client {
	return &Client{BaseUrl: "https://api.modrinth.com/v2"}
}
