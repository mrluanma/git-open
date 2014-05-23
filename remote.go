package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
)

type Remote struct {
	Name, Url, Type string
}

func DetectRemote(dir string) ([]*Remote, error) {
	git := exec.Command("git", "remote", "-v")
	git.Dir = dir

	stdout, err := git.StdoutPipe()
	if err != nil {
		return nil, err
	}
	git.Stderr = os.Stderr

	err = git.Start()
	if err != nil {
		return nil, err
	}

	remotes := make([]*Remote, 0)

	reader := bufio.NewReader(stdout)
	re := regexp.MustCompile(`^(\S+)\s*(\S+)\s*\((\S+)\)`)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		m := re.FindStringSubmatch(line)
		if m != nil {
			r := &Remote{m[1], m[2], m[3]}
			remotes = append(remotes, r)
		}
	}

	if err := git.Wait(); err != nil {
		return nil, err
	}

	return remotes, nil
}

// Convert git remote url to github/bitbucket URL
//
// Supported types:
//   - git@github.com:username/repo.git
//   - git@github.com:username/repo
//   - http://github.com/username/repo.git
//   - http://github.com/username/repo
//   - https://github.com/username/repo.git
//   - https://github.com/username/repo
//   - git://github.com/username/repo.git
//   - git://github.com/username/repo
func MangleURL(url string) (string, error) {
	ssh_re := regexp.MustCompile(`^git@(.*?):(.*?)/(.*?)(\.git)?$`)
	https_re := regexp.MustCompile(`^https?://(.*?)/(.*?)/(.*?)(\.git)?$`)
	git_re := regexp.MustCompile(`^git://(.*?)/(.*?)/(.*?)(\.git)?$`)

	var matches []string

	if m := ssh_re.FindStringSubmatch(url); m != nil {
		matches = m
	} else if m := https_re.FindStringSubmatch(url); m != nil {
		matches = m
	} else if m := git_re.FindStringSubmatch(url); m != nil {
		matches = m
	} else {
		return "", fmt.Errorf("unsupported remote url: %s", url)
	}

	return CreateURL(matches[1], matches[2], matches[3])
}

func CreateURL(host, user, repo string) (string, error) {
	return fmt.Sprintf("https://%s/%s/%s", host, user, repo), nil
}
