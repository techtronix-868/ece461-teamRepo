package datastore

import (
	"log"
	"strconv"

	"github.com/mabaums/ece461-web/backend/models"
)

type InMemoryDatstore struct {
	packages    []models.PackageMetadata
	initialized bool
}

func (md *InMemoryDatstore) initIfEmpty() {
	if !md.initialized {
		log.Printf("Initializing InMemoryDatastore")
		md.initialized = true
		for i := 0; i <= 20; i++ {
			pkg := models.PackageMetadata{
				Name:    "Package",
				Version: "1.0.0",
				ID:      strconv.Itoa(i),
			}
			md.packages = append(md.packages, pkg)
		}
		// todo: add some info about packages here
	}
}

func (md *InMemoryDatstore) GetPackage(id string) (*models.Package, *models.Error) {
	md.initIfEmpty()
	return nil, nil // return actual package
}
func (md *InMemoryDatstore) GetPackages() []models.PackageMetadata {
	md.initIfEmpty()
	return md.packages
}
