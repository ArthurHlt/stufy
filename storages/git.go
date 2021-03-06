package storages

import (
	"fmt"
	"github.com/ArthurHlt/stufy/loading"
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
	loading.Start(fmt.Sprintf("Cloning %s in temp dir ", target), "", nil)
	defer loading.Stop()
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:          target,
		Auth:         authMethod,
		SingleBranch: true,
	})
	if err != nil {
		panic(err)
	}
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
	if !messages.StopShow() {
		spin.Start()
	}
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

func (s Git) Resync() error {
	os.RemoveAll(s.folder)
	loading.Start(
		fmt.Sprintf("Resynchronize %s in temp dir ", s.target),
		fmt.Sprintf("Finished resynchronize %s\n", s.target),
		nil,
	)
	defer loading.Stop()
	_, err := git.PlainClone(s.folder, false, &git.CloneOptions{
		URL:          s.target,
		Auth:         s.authMethod,
		SingleBranch: true,
	})
	return err
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
	loading.Start(
		"Push changes in remote",
		"Finished push changes in remote",
		nil,
	)
	defer func() {
		loading.Stop()
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
