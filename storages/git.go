package storages

import (
	"fmt"
	"github.com/ArthurHlt/stufy/messages"
	"github.com/ArthurHlt/stufy/model"
	"github.com/briandowns/spinner"
	"github.com/mitchellh/go-homedir"
	"github.com/whilp/git-urls"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Git struct {
	folder     string
	target     string
	local      *Local
	repo       *git.Repository
	authMethod transport.AuthMethod
	pulled     bool
}

func NewGit(target string) *Git {
	tmpDir := os.TempDir()
	uri, err := giturls.Parse(target)
	if err != nil {
		panic(err)
	}
	splitPath := strings.Split(strings.TrimPrefix(uri.Path, "/"), "/")
	dir := filepath.Join(tmpDir, "stufy", uri.Host)
	for _, p := range splitPath {
		dir = filepath.Join(dir, p)
	}

	authMethod, err := createAuthMethod(uri)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(dir)

	if err == nil {
		repo, err := git.PlainOpen(dir)
		if err != nil {
			panic(err)
		}
		return &Git{
			folder:     dir,
			target:     target,
			repo:       repo,
			authMethod: authMethod,
			local:      NewLocal(dir),
		}
	}
	spin := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	spin.Prefix = fmt.Sprintf("Cloning %s in temp dir ", target)
	spin.Start()
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:          target,
		Auth:         authMethod,
		SingleBranch: true,
	})
	if err != nil {
		spin.Stop()
		panic(err)
	}
	spin.Stop()
	return &Git{
		folder:     dir,
		target:     target,
		repo:       repo,
		authMethod: authMethod,
		local:      NewLocal(dir),
		pulled:     true,
	}
}

func (s Git) Config() (model.Config, error) {
	return s.local.Config()
}

func (s Git) Incidents() (model.Incidents, error) {
	if s.pulled {
		return s.local.Incidents()
	}
	w, err := s.repo.Worktree()
	if err != nil {
		return model.Incidents{}, err
	}
	spin := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	spin.Prefix = "Pulling changes from remote "
	spin.Start()
	err = w.Pull(&git.PullOptions{
		Auth:         s.authMethod,
		RemoteName:   "origin",
		SingleBranch: true,
	})
	if err != nil && err.Error() != git.NoErrAlreadyUpToDate.Error() {
		spin.Stop()
		return model.Incidents{}, err
	}
	spin.Stop()
	return s.local.Incidents()
}

func (s Git) CreateIncident(incident model.Incident) error {
	err := s.local.CreateIncident(incident)
	if err != nil {
		return err
	}
	return s.pushCommit(incident, "Add incident", false)
}

func (s Git) UpdateIncident(incident model.Incident) error {
	err := s.local.UpdateIncident(incident)
	if err != nil {
		return err
	}
	return s.pushCommit(incident, "Update incident", false)
}

func (s Git) DeleteIncident(incident model.Incident) error {
	err := s.local.DeleteIncident(incident)
	if err != nil {
		return err
	}
	return s.pushCommit(incident, "Delete incident", true)
}

func (s Git) Open(incident model.Incident) error {
	err := s.local.Open(incident)
	if err != nil {
		return err
	}
	return s.pushCommit(incident, "Update incident", false)
}

func (Git) Detect(target string) bool {
	return strings.HasSuffix(target, ".git") || strings.HasPrefix(target, "git://")
}

func (s Git) pushCommit(incident model.Incident, message string, remove bool) error {
	w, err := s.repo.Worktree()
	if err != nil {
		return err
	}
	path := filepath.Join(incidentFolder, incident.Filename())
	if !remove {
		_, err = w.Add(path)
	} else {
		_, err = w.Remove(path)
	}

	if err != nil {
		return err
	}
	_, err = w.Commit(fmt.Sprintf("%s %s", message, incident.Filename()), &git.CommitOptions{
		Author: signature(),
	})
	if err != nil {
		return err
	}
	spin := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	spin.Prefix = "Push changes in remote"
	spin.FinalMSG = "Finished push changes in remote"
	spin.Start()
	defer func() {
		spin.Stop()
		messages.Print("\r\n")
	}()
	err = s.repo.Push(&git.PushOptions{
		Auth:       s.authMethod,
		RemoteName: "origin",
	})
	if err != nil && err.Error() != git.NoErrAlreadyUpToDate.Error() {
		err = w.Pull(&git.PullOptions{
			Auth:         s.authMethod,
			RemoteName:   "origin",
			SingleBranch: true,
		})
		if err != nil {
			return err
		}
	}
	err = s.repo.Push(&git.PushOptions{
		Auth:       s.authMethod,
		RemoteName: "origin",
	})
	if err.Error() != git.NoErrAlreadyUpToDate.Error() {
		return err
	}
	return nil
}

func createAuthMethod(uri *url.URL) (transport.AuthMethod, error) {
	if uri.Scheme == "ssh" {
		return ssh.NewSSHAgentAuth(uri.User.Username())
	}
	if uri.User == nil || uri.User.Username() == "" {
		return nil, nil
	}
	password, _ := uri.User.Password()
	return &http.BasicAuth{
		Username: uri.User.Username(),
		Password: password,
	}, nil
}

func signature() *object.Signature {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(home, ".gitconfig"))
	if err != nil {
		return &object.Signature{
			Name:  "stufy",
			Email: "stufy@github.com",
			When:  time.Now(),
		}
	}
	var confGit config.Config
	d := config.NewDecoder(f)
	err = d.Decode(&confGit)
	if err != nil {
		panic(err)
	}
	section := confGit.Section("user")

	return &object.Signature{
		Name:  section.Option("name"),
		Email: section.Option("email"),
		When:  time.Now(),
	}
}
