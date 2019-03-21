package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	MajorOutage         Severity = "major-outage"
	PartialOutage       Severity = "partial-outage"
	DegradedPerformance Severity = "degraded-performance"
	UnderMaintenance    Severity = "under-maintenance"
)

var defaultSeverities = []Severity{
	MajorOutage,
	PartialOutage,
	DegradedPerformance,
	UnderMaintenance,
}

func FindSeverity(severity string) (Severity, error) {
	for _, s := range defaultSeverities {
		if s == Severity(severity) {
			return s, nil
		}
	}
	return "", fmt.Errorf("Sevrity '%s' is not valid", severity)
}

func Severities() []Severity {
	return defaultSeverities
}

func SeveritiesString() []string {
	severities := make([]string, len(defaultSeverities))
	for i, s := range defaultSeverities {
		severities[i] = string(s)
	}
	return severities
}

type Severity string

type Incidents []Incident

func (is Incidents) Filenames() []string {
	selects := make([]string, len(is))
	for i, incident := range is {
		selects[i] = incident.Filename()
	}
	return selects
}

type Incident struct {
	Id              string   `yaml:"id,omitempty" json:"id,omitempty"`
	Title           string   `yaml:"title" json:"title,omitempty"`
	Description     string   `yaml:"description,omitempty" json:"description,omitempty"`
	Date            Date     `yaml:"date" json:"date,omitempty"`
	Modified        Date     `yaml:"modified,omitempty" json:"modified,omitempty"`
	Severity        Severity `yaml:"severity,omitempty" json:"severity,omitempty"`
	AffectedSystems []string `yaml:"affectedsystems" json:"affectedsystems,omitempty"`
	Resolved        bool     `yaml:"resolved" json:"resolved,omitempty"`
	Scheduled       *Date    `yaml:"scheduled,omitempty" json:"scheduled,omitempty"`
	Duration        int      `yaml:"duration,omitempty" json:"duration,omitempty"`
	Content         string   `yaml:"-" fm:"content"`
}

func (i Incident) String() string {
	return i.Filename()
}

func (i Incident) Filename() string {
	t := time.Time(i.Date).Format("2006-01-02")
	title := strings.ToLower(strings.Replace(i.Title, " ", "_", -1))
	return fmt.Sprintf("%s_%s.md", t, title)
}

type Date time.Time

func (d Date) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d *Date) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t = t.In(time.Now().Location())
	*d = Date(t)

	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)

	return json.Marshal(t.Format(time.RFC3339))
}

func (d Date) MarshalYAML() (interface{}, error) {
	t := time.Time(d)
	return t.UTC().Format(time.RFC3339), nil
}
