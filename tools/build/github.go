package build

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubAccount struct {
	Name  string
	Email string
	Org   string
	Repo  string
	Token string
}

func newGithubClient(account *GithubAccount) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: account.Token},
	)
	return github.NewClient(oauth2.NewClient(ctx, ts)), nil
}

func GithubCreateRelease(ctx context.Context, account *GithubAccount, tag string, pre bool) (int64, error) {
	client, err := newGithubClient(account)
	if err != nil {
		return 0, err
	}
	release, _, err := client.Repositories.CreateRelease(ctx, account.Org, account.Repo, &github.RepositoryRelease{
		TagName:    &tag,
		Prerelease: &pre,
	})
	if err != nil {
		return 0, err
	}
	return release.GetID(), nil
}

func GithubUploadAsset(ctx context.Context, account *GithubAccount, id int64, file string) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	client, err := newGithubClient(account)
	if err != nil {
		return err
	}
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, account.Org, account.Repo, id, &github.UploadOptions{
		Name: filepath.Base(file),
	}, fileReader)
	return err
}

func GithubGetFileSHA(ctx context.Context, account *GithubAccount, path string) (string, error) {
	client, err := newGithubClient(account)
	if err != nil {
		return "", err
	}
	content, _, _, err := client.Repositories.GetContents(ctx, account.Org, account.Repo, path, nil)
	if err != nil {
		return "", err
	}
	return content.GetSHA(), nil
}

func GithubUploadFile(ctx context.Context, account *GithubAccount, path string, file string) error {
	client, err := newGithubClient(account)
	if err != nil {
		return err
	}
	sha, err := GithubGetFileSHA(ctx, account, path)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	message := "update"
	_, _, err = client.Repositories.UpdateFile(ctx, account.Org, account.Repo, path, &github.RepositoryContentFileOptions{
		Content: content,
		SHA:     &sha,
		Message: &message,
		Committer: &github.CommitAuthor{
			Name:  &account.Name,
			Email: &account.Email,
		},
	})
	return err
}
