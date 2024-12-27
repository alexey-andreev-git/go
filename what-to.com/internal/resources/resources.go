package resources

import (
	"embed"
)

var (
	//go:embed appfs
	appResources embed.FS
)

type (
	AppSourcesInterface interface {
		GetFs() embed.FS
	}
	AppSources struct{}
)

func NewAppSources() *AppSources {
	return &AppSources{}
}

func (us *AppSources) GetFs() embed.FS {
	return appResources
}
