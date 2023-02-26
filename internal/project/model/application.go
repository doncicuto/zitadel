package model

import (
	"net/url"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Application struct {
	es_models.ObjectRoot

	AppID      string
	State      AppState
	Name       string
	ExternalURL string
	IsVisibleToEndUser bool
	Type       AppType
	OIDCConfig *OIDCConfig
	APIConfig  *APIConfig
	SAMLConfig *SAMLConfig
}

type AppState int32

const (
	AppStateActive AppState = iota
	AppStateInactive
	AppStateRemoved
)

type AppType int32

const (
	AppTypeUnspecified AppType = iota
	AppTypeOIDC
	AppTypeSAML
	AppTypeAPI
)

func (a *Application) IsValid(includeConfig bool) bool {
	if a.Name == "" || a.AggregateID == "" {
		return false
	}
	if a.ExternalURL != "" {
		u, err := url.Parse(a.ExternalURL)
		if err != nil || u.Scheme == "" && u.Host == "" {			
			return false
		}
	}
	if !includeConfig {
		return true
	}
	if a.Type == AppTypeOIDC && !a.OIDCConfig.IsValid() {
		return false
	}
	if a.Type == AppTypeAPI && !a.APIConfig.IsValid() {
		return false
	}
	if a.Type == AppTypeSAML && !a.SAMLConfig.IsValid() {
		return false
	}
	return true
}
