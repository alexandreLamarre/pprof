package ui

import (
	"github.com/google/pprof/internal/driver"
	"github.com/google/pprof/profile"
)

func NewWebUI(p *profile.Profile) (*driver.WebUI, error) {
	return driver.NewWebUI(p)
}
