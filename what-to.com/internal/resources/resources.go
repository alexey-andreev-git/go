package resources

import (
	"embed"
)

var (
	//go:embed appfs
	appResources embed.FS
)

type AppSourcesInterface interface {
	GetRes() embed.FS
}

type AppSources struct{}

func NewAppSources() AppSourcesInterface {
	return &AppSources{}
}

func (us *AppSources) GetRes() embed.FS {
	return appResources
}
