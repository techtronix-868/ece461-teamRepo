package datastore

import "github.com/mabaums/ece461-web/backend/models"

type InMemoryDatstore struct {
	packages []models.PackageMetadata
}

func (md InMemoryDatstore) GetPackages() []models.PackageMetadata {
	return md.packages
}
