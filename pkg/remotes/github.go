package remotes

import (
	"context"
	"net/http"

	"github.com/deps-cloud/discovery/pkg/config"

	"github.com/google/go-github/v20/github"

	"github.com/sirupsen/logrus"

	"golang.org/x/oauth2"
)

// NewGithubRemote constructs a new remote implementation that speaks with GitHub
// for repository related information.
func NewGithubRemote(cfg *config.Github) (Remote, error) {
	baseURL := cfg.GetBaseUrl()
	uploadURL := cfg.GetUploadUrl()

	fn := func(client *http.Client) (*github.Client, error) {
		return github.NewClient(client), nil
	}
	if baseURL != nil && uploadURL != nil {
		fn = func(client *http.Client) (*github.Client, error) {
			return github.NewEnterpriseClient(baseURL.Value, uploadURL.Value, client)
		}
	}

	httpClient := &http.Client{}
	if o2 := cfg.GetOauth2(); o2 != nil {
		ts := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: o2.Token,
		})

		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	client, err := fn(httpClient)
	if err != nil {
		return nil, err
	}

	return &githubRemote{
		config: cfg,
		client: client,
	}, nil
}

var _ Remote = &githubRemote{}

type githubRemote struct {
	config *config.Github
	client *github.Client
}

func (r *githubRemote) ListRepositories() ([]string, error) {
	organizations := make([]string, 0)
	repositories := make([]string, 0)

	if r.config.Organizations != nil {
		// init with configured orgs
		organizations = append(organizations, r.config.Organizations...)
	}

	// discover more from users
	for _, user := range r.config.Users {
		logrus.Infof("[remotes.github] processing organizations for user: %s", user)

		for orgPage := 1; orgPage != 0; {
			orgs, response, err := r.client.Organizations.List(context.Background(), user, &github.ListOptions{
				Page: orgPage,
			})

			if err != nil {
				logrus.Errorf("[remotes.github] encountered err on orgPage %d, %v", orgPage, err)
				break
			}

			orgLogins := make([]string, len(orgs))
			for i, org := range orgs {
				orgLogins[i] = org.GetLogin()
			}

			organizations = append(organizations, orgLogins...)

			orgPage = response.NextPage
		}

		logrus.Infof("[remotes.github] processing repositories for user: %s", user)

		for repoPage := 1; repoPage != 0; {
			repos, response, err := r.client.Repositories.List(context.Background(), user, &github.RepositoryListOptions{
				ListOptions: github.ListOptions{
					Page: repoPage,
				},
			})

			if err != nil {
				logrus.Errorf("[remotes.github] encountered err on repoPage %d, %v", repoPage, err)
				break
			}

			urls := make([]string, len(repos))

			for i, repo := range repos {
				if r.config.GetStrategy() == config.CloneStrategy_HTTP {
					urls[i] = repo.GetCloneURL()
				} else {
					urls[i] = repo.GetSSHURL()
				}
			}

			repositories = append(repositories, urls...)

			repoPage = response.NextPage
		}
	}

	for _, organization := range organizations {
		logrus.Infof("[remotes.github] processing repositories for organization: %s", organization)

		for orgRepoPage := 1; orgRepoPage != 0; {
			orgRepos, response, err := r.client.Repositories.ListByOrg(context.Background(), organization, &github.RepositoryListByOrgOptions{
				ListOptions: github.ListOptions{
					Page: orgRepoPage,
				},
			})

			if err != nil {
				logrus.Errorf("[remotes.github] encountered err on orgRepoPage %d, %v", orgRepoPage, err)
				break
			}

			urls := make([]string, len(orgRepos))

			for i, orgRepo := range orgRepos {
				if r.config.GetStrategy() == config.CloneStrategy_HTTP {
					urls[i] = orgRepo.GetCloneURL()
				} else {
					urls[i] = orgRepo.GetSSHURL()
				}
			}

			repositories = append(repositories, urls...)

			orgRepoPage = response.NextPage
		}
	}

	return repositories, nil
}
