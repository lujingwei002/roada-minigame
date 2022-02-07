package cli

var defaultClient *Client

func Default() *Client {
	return defaultClient
}

func init() {
	defaultClient = NewClient()
}
