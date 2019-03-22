package storages

import (
	"github.com/ArthurHlt/stufy/model"
	"os"
	"strings"
)

const (
	configFilename = "config.yml"
	incidentFolder = "content"
)

type Storage interface {
	Config() (model.Config, error)
	Incidents() (model.Incidents, error)
	CreateIncident(model.Incident) error
	UpdateIncident(model.Incident) error
	DeleteIncident(model.Incident) error
	Open(model.Incident) error
	Resync() error
}

func FindStorage(target string) Storage {
	if strings.HasSuffix(target, ".git") || strings.HasPrefix(target, "git://") {
		return NewGit(target)
	}
	if _, err := os.Stat(target); err == nil {
		return NewLocal(target)
	}
	return nil
}
