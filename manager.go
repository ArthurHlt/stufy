package stufy

import (
	"fmt"
	"github.com/ArthurHlt/stufy/model"
	"github.com/ArthurHlt/stufy/storages"
	"os"
	"strings"
	"time"
)

type Manager struct {
	storage        storages.Storage
	incidentsCache model.Incidents
	dotConfig      *DotConfig
	target         string
}

func NewManager(target string) (*Manager, error) {
	dotConfig, err := loadDotConfig()
	if err != nil {
		return nil, err
	}
	dotTarget := dotConfig.Aliases.FindTarget(target)
	if dotTarget != "" {
		target = dotTarget
	}
	if target == "" {
		target, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	s := storages.FindStorage(target)
	if s == nil {
		return nil, fmt.Errorf("Can't find storage for '%s'.", target)
	}

	return &Manager{
		storage:        s,
		incidentsCache: make(model.Incidents, 0),
		dotConfig:      dotConfig,
		target:         target,
	}, nil
}

func (m Manager) Config() (model.Config, error) {
	return m.storage.Config()
}

func (m Manager) Aliases() Aliases {
	return m.dotConfig.Aliases
}

func (m Manager) AddAlias(alias string) error {
	m.dotConfig.Aliases = append(m.dotConfig.Aliases, Alias{
		Name:   alias,
		Target: m.target,
	})
	return saveDotConfig(m.dotConfig)
}

func (m Manager) RemoveAlias(alias string) error {
	newAliases := make([]Alias, 0)
	for _, a := range m.dotConfig.Aliases {
		if a.Name == alias {
			continue
		}
		newAliases = append(newAliases, a)
	}
	m.dotConfig.Aliases = newAliases
	return saveDotConfig(m.dotConfig)
}

func (m Manager) CreateIncident(req RequestCreate) error {
	severity, err := model.FindSeverity(req.Severity)
	if err != nil {
		return err
	}
	incident := model.Incident{
		Title:           req.Cause,
		Date:            model.Date(time.Now()),
		Content:         req.Description,
		AffectedSystems: req.Systems,
		Severity:        severity,
	}
	if req.Open {
		incident, err = m.storage.Open(incident)
		if err != nil {
			return err
		}
	}
	err = m.storage.CreateIncident(incident)
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) DeleteIncident(req RequestDelete) error {
	if !req.Confirm {
		return nil
	}
	incident, err := m.FindIncident(req.Filename)
	if err != nil {
		return err
	}
	return m.storage.DeleteIncident(incident)
}

func (m Manager) UpdateIncident(req RequestUpdate) error {
	if !req.Confirm {
		return nil
	}
	incident, err := m.FindIncident(req.Filename)
	if err != nil {
		return err
	}

	mod := model.Date(time.Now())
	content := incident.Content

	if req.UpdateType != "" && req.UpdateType != "no" {
		content += fmt.Sprintf(
			"\n\n::: update %s | %s\n%s\n:::",
			strings.Title(req.UpdateType),
			time.Time(mod).UTC().Format(time.RFC3339),
			req.UpdateContent,
		)
	}
	if req.Severity != "" {
		severity, err := model.FindSeverity(req.Severity)
		if err != nil {
			return err
		}
		incident.Severity = severity
	}

	incident.Content = content
	if len(req.Systems) != 0 {
		incident.AffectedSystems = req.Systems
	}
	incident.Modified = mod
	incident.Resolved = req.Resolved

	if req.Open {
		incident, err = m.storage.Open(incident)
		if err != nil {
			return err
		}
	}
	err = m.storage.UpdateIncident(incident)
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) CreateScheduled(req RequestScheduled) error {
	dur, err := time.ParseDuration(req.Duration)
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02T15:04", req.Date, time.Local)
	if err != nil {
		return err
	}
	t = t.UTC()
	d := model.Date(t)
	incident := model.Incident{
		Title:           req.Title,
		Date:            model.Date(time.Now()),
		Content:         req.Description,
		AffectedSystems: req.Systems,
		Scheduled:       &d,
		Duration:        int(dur.Minutes()),
	}
	if req.Open {
		incident, err = m.storage.Open(incident)
		if err != nil {
			return err
		}
	}
	err = m.storage.CreateIncident(incident)
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) UpdateScheduled(req RequestUpdateScheduled) error {
	if !req.Confirm {
		return nil
	}

	incident, err := m.FindIncident(req.Filename)
	if err != nil {
		return err
	}
	if req.Description != "" {
		incident.Content = req.Description
	}

	if req.Date != "" {
		t, err := time.ParseInLocation("2006-01-02T15:04", req.Date, time.Local)
		if err != nil {
			return err
		}
		t = t.UTC()
		d := model.Date(t)
		incident.Scheduled = &d
	}

	if len(req.Systems) != 0 {
		incident.AffectedSystems = req.Systems
	}

	if req.Duration != "" {
		dur, err := time.ParseDuration(req.Duration)
		if err != nil {
			return err
		}
		incident.Duration = int(dur.Minutes())
	}

	if req.Open {
		incident, err = m.storage.Open(incident)
		if err != nil {
			return err
		}
	}
	err = m.storage.UpdateIncident(incident)
	if err != nil {
		return err
	}
	return nil
}

func (m Manager) FinishScheduled(req RequestUnscheduled) error {
	if !req.Confirm {
		return nil
	}
	incident, err := m.FindIncident(req.Filename)
	if err != nil {
		return err
	}
	incident.Resolved = true
	return m.storage.UpdateIncident(incident)
}

func (m Manager) DeleteScheduled(req RequestUnscheduled) error {
	if !req.Confirm {
		return nil
	}
	incident, err := m.FindIncident(req.Filename)
	if err != nil {
		return err
	}
	return m.storage.DeleteIncident(incident)
}

func (m Manager) FindIncident(incidentFilename string) (model.Incident, error) {
	err := m.cacheIncidents()
	if err != nil {
		return model.Incident{}, err
	}
	for _, i := range m.incidentsCache {
		if i.Filename() == incidentFilename {
			return i, nil
		}
	}
	return model.Incident{}, fmt.Errorf("Can't find incident '%s'", incidentFilename)
}

func (m *Manager) cacheIncidents() error {
	if len(m.incidentsCache) > 0 {
		return nil
	}
	incidents, err := m.storage.Incidents()
	if err != nil {
		return err
	}
	m.incidentsCache = incidents
	return nil
}

func (m *Manager) ListIncident(showResolved bool) (model.Incidents, error) {
	filtered := make([]model.Incident, 0)
	var incidents []model.Incident
	var err error
	if len(m.incidentsCache) > 0 {
		incidents = m.incidentsCache
	} else {
		incidents, err = m.storage.Incidents()
		if err != nil {
			return filtered, err
		}
		m.incidentsCache = incidents
	}

	for _, i := range incidents {
		if i.Scheduled != nil {
			continue
		}
		if i.Resolved && !showResolved {
			continue
		}
		filtered = append(filtered, i)
	}
	return filtered, nil
}

func (m *Manager) Resync() error {
	return m.storage.Resync()
}

func (m *Manager) ListScheduled(showResolved bool) (model.Incidents, error) {
	filtered := make([]model.Incident, 0)
	var incidents []model.Incident
	var err error
	if len(m.incidentsCache) > 0 {
		incidents = m.incidentsCache
	} else {
		incidents, err = m.storage.Incidents()
		if err != nil {
			return filtered, err
		}
	}

	for _, i := range incidents {
		if i.Scheduled == nil {
			continue
		}
		if i.Resolved && !showResolved {
			continue
		}
		filtered = append(filtered, i)
	}
	return filtered, nil
}
