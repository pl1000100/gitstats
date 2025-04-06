package github

type APIClient struct {
	Token string
}

func NewGitHubAPIClient(token string) *APIClient {
	return &APIClient{
		Token: token,
	}
}
