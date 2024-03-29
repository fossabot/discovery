package config_test

import (
	"testing"

	"github.com/deps-cloud/discovery/pkg/config"

	"github.com/stretchr/testify/require"
)

func testBasic(t *testing.T, basic *config.Basic) {
	require.NotNil(t, basic)
	require.Equal(t, "username", basic.Username)

	require.NotNil(t, basic.Password)
	require.Equal(t, "password", basic.Password.Value)
}

func testOauth(t *testing.T, token *config.OAuthToken) {
	require.NotNil(t, token)
	require.Equal(t, "token", token.Token)

	if token.ApplicationId != nil {
		require.Equal(t, "application_id", token.ApplicationId.Value)
	}
}

func testOauth2(t *testing.T, token *config.OAuth2Token) {
	require.NotNil(t, token)
	require.Equal(t, "token", token.Token)

	require.NotNil(t, token.TokenType)
	require.Equal(t, "token_type", token.TokenType.Value)

	require.NotNil(t, token.RefreshToken)
	require.Equal(t, "refresh_token", token.RefreshToken.Value)

	require.NotNil(t, token.Expiry)
	require.Equal(t, "expiry", token.Expiry.Value)
}

func testGeneric(t *testing.T, generic *config.Generic) {
	require.NotNil(t, generic)
	require.Equal(t, "base_url", generic.BaseUrl)
	require.Equal(t, "path", generic.Path)
	require.Equal(t, "per_page_parameter", generic.PerPageParameter)
	require.Equal(t, "page_parameter", generic.PageParameter)
	require.Equal(t, int32(20), generic.PageSize)
	require.Equal(t, "selector", generic.Selector)
}

func testGitlab(t *testing.T, gitlab *config.Gitlab) {
	require.NotNil(t, gitlab)
	require.NotNil(t, gitlab.BaseUrl)

	require.Equal(t, "base_url", gitlab.BaseUrl.Value)
	require.Equal(t, config.CloneStrategy_HTTP, gitlab.Strategy)
}

func testGithub(t *testing.T, github *config.Github) {
	require.NotNil(t, github)
	require.NotNil(t, github.BaseUrl)
	require.Equal(t, "base_url", github.BaseUrl.Value)

	require.NotNil(t, github.UploadUrl)
	require.Equal(t, "upload_url", github.UploadUrl.Value)

	require.NotNil(t, github.Organizations)
	require.Len(t, github.Organizations, 1)
	require.Contains(t, github.Organizations, "org1")

	require.NotNil(t, github.Users)
	require.Len(t, github.Users, 1)
	require.Contains(t, github.Users, "user1")

	require.Equal(t, config.CloneStrategy_HTTP, github.Strategy)
}

func testBitbucket(t *testing.T, bitbucket *config.Bitbucket) {
	require.NotNil(t, bitbucket)

	require.NotNil(t, bitbucket.Teams)
	require.Len(t, bitbucket.Teams, 1)
	require.Contains(t, bitbucket.Teams, "team1")

	require.NotNil(t, bitbucket.Users)
	require.Len(t, bitbucket.Users, 1)
	require.Contains(t, bitbucket.Users, "user1")

	require.Equal(t, config.CloneStrategy_HTTP, bitbucket.Strategy)
}

func testStatic(t *testing.T, static *config.Static) {
	require.NotNil(t, static)
	require.Len(t, static.RepositoryUrls, 1)

	require.Contains(t, static.RepositoryUrls, "repository_urls")
}

func testRds(t *testing.T, rds *config.Rds) {
	require.NotNil(t, rds)
	require.Equal(t, rds.Target, "target")
}

func testCommon(t *testing.T, cfg *config.Configuration) {
	require.Len(t, cfg.Accounts, 10)

	{
		generic := cfg.Accounts[0].GetGeneric()
		testGeneric(t, generic)
	}

	{
		generic := cfg.Accounts[1].GetGeneric()
		testGeneric(t, generic)
		testBasic(t, generic.Basic)
	}

	{
		gitlab := cfg.Accounts[2].GetGitlab()
		testGitlab(t, gitlab)
		testOauth(t, gitlab.Private)
	}

	{
		gitlab := cfg.Accounts[3].GetGitlab()
		testGitlab(t, gitlab)
		testOauth(t, gitlab.Oauth)
	}

	{
		github := cfg.Accounts[4].GetGithub()
		testGithub(t, github)
	}

	{
		github := cfg.Accounts[5].GetGithub()
		testGithub(t, github)
		testOauth2(t, github.Oauth2)
	}

	{
		bitbucket := cfg.Accounts[6].GetBitbucket()
		testBitbucket(t, bitbucket)
		testBasic(t, bitbucket.Basic)
	}

	{
		bitbucket := cfg.Accounts[7].GetBitbucket()
		testBitbucket(t, bitbucket)
		testOauth(t, bitbucket.Oauth)
	}

	{
		static := cfg.Accounts[8].GetStatic()
		testStatic(t, static)
	}

	{
		rds := cfg.Accounts[9].GetRds()
		testRds(t, rds);
	}
}

func Test_proto(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.prototxt")
	require.NoError(t, err)
	testCommon(t, cfg)
}

func Test_yaml(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.yaml")
	require.NoError(t, err)
	testCommon(t, cfg)
}

func Test_json(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.json")
	require.NoError(t, err)
	testCommon(t, cfg)
}
