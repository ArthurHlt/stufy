package storages

import (
	"fmt"
	"github.com/ArthurHlt/open-golang/open"
	"github.com/ArthurHlt/stufy/model"
	"github.com/ericaro/frontmatter"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Local struct {
	folder string
}

func NewLocal(folder string) *Local {
	return &Local{
		folder: folder,
	}
}

func (l Local) Config() (model.Config, error) {
	confPath := filepath.Join(l.folder, configFilename)
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		return model.Config{}, err
	}
	var config model.Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}

func (l Local) Incidents() (model.Incidents, error) {
	dir := filepath.Join(l.folder, incidentFolder)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	incidents := make([]model.Incident, 0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		tmpIncidentPath := filepath.Join(dir, file.Name())
		if filepath.Ext(tmpIncidentPath) != ".md" {
			continue
		}
		b, err := ioutil.ReadFile(tmpIncidentPath)
		if err != nil {
			return incidents, err
		}
		var tmpIncident model.Incident
		err = frontmatter.Unmarshal(b, &tmpIncident)
		if err != nil {
			return incidents, err
		}
		incidents = append(incidents, tmpIncident)
	}
	return incidents, nil
}

func (l Local) Resync() error {
	return fmt.Errorf("no resync is available on local target")
}

func (l Local) CreateIncident(incident model.Incident) error {
	b, err := frontmatter.Marshal(incident)
	if err != nil {
		return nil
	}
	incidentFile := filepath.Join(l.folder, incidentFolder, incident.Filename())
	return ioutil.WriteFile(incidentFile, b, 0655)
}

func (l Local) UpdateIncident(incident model.Incident) error {
	return l.CreateIncident(incident)
}

func (l Local) DeleteIncident(incident model.Incident) error {
	incidentFile := filepath.Join(l.folder, incidentFolder, incident.Filename())
	return os.Remove(incidentFile)
}

func (l Local) Open(incident model.Incident) error {
	editor := ""
	if runtime.GOOS == "windows" {
		editor = "notepad"
	}
	if v := os.Getenv("VISUAL"); v != "" {
		editor = v
	} else if e := os.Getenv("EDITOR"); e != "" {
		editor = e
	}
	incidentFile := filepath.Join(l.folder, incidentFolder, incident.Filename())
	if editor != "" {
		return open.RunWithWait(incidentFile, editor)
	}

	return open.RunWait(incidentFile)
}
