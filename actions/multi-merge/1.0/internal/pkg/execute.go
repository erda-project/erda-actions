package pkg

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	toolScript = "/opt/action/common.sh"
)

func (m *MultiMerge) Execute() error {
	defer func() {
		if err := m.results.Store(); err != nil {
			logrus.Errorf("failed to store results: %v", err)
		}
		return
	}()
	// initialize git
	initGitCmd := NewGitInitCommand("load_pubkey", "load_git_crypt_key", "configure_https_tunnel", "configure_git_ssl_verification", "configure_credentials")
	if err := initGitCmd.Execute(); err != nil {
		m.results.Add("init-git", "failed")
		return err
	}
	if m.cfg.Username != "" && m.cfg.Password != "" {
		loginCmd := NewGitLoginCommand(m.cfg.Username, m.cfg.Password)
		if err := loginCmd.Execute(); err != nil {
			m.results.Add("login", "failed")
			return err
		}
	}
	if len(m.cfg.GitConfigs) > 0 {
		configCmd := NewGitConfigCommand(m.cfg.GitConfigs)
		if err := configCmd.Execute(); err != nil {
			m.results.Add("config", "failed")
			return err
		}
	}

	// clone the dest repo
	cloneCmd := NewGitCloneCommand(m.cfg.DestRepo, m.cfg.DestBranch)
	if err := cloneCmd.Execute(); err != nil {
		return err
	}
	m.results.Add("clone-dest", "success")

	// merge all repos
	mergedRepos := make([]string, 0)
	for _, repo := range m.cfg.Repos {
		mergeCmd := NewGitMergeCommand(repo.Uri, repo.Branches)
		if err := mergeCmd.Execute(); err != nil {
			m.results.Add("merged-repos", fmt.Sprintf("%v", mergedRepos))
			return err
		}
		mergedRepos = append(mergedRepos, repo.Uri)
	}
	m.results.Add("merged-repos", fmt.Sprintf("%v", mergedRepos))
	return nil
}
