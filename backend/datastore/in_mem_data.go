package datastore

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/mabaums/ece461-web/backend/models"
)

type InMemoryDatstore struct {
	packages    []models.Package
	initialized bool
}

func (md *InMemoryDatstore) initIfEmpty() {
	if !md.initialized {
		log.Info("Initializing InMemoryDatastore")
		md.initialized = true
		for i := 0; i <= 20; i++ {
			pkg_meta := models.PackageMetadata{
				Name:    "Package",
				Version: "1.0.0",
				ID:      strconv.Itoa(i),
			}
			pkg_data := models.PackageData{
				Content:   "Content",
				URL:       "www.google.com",
				JSProgram: "string",
			}
			pkg := models.Package{Data: pkg_data, Metadata: pkg_meta}
			md.packages = append(md.packages, pkg)
		}
	}
}

func (md *InMemoryDatstore) GetPackage(id string) (*models.Package, *models.Error) {
	md.initIfEmpty()
	log.Infof("Getting Package id: %v\n", id)
	for _, pkg := range md.packages {
		if pkg.Metadata.ID == id {
			return &pkg, nil
		}
	}
	return nil, &models.Error{Code: http.StatusNotFound, Message: "Package with id not found"} // return actual package
}
func (md *InMemoryDatstore) GetPackages() []models.PackageMetadata {
	md.initIfEmpty()
	met_pkgs := []models.PackageMetadata{}
	for _, pkg := range md.packages {
		met_pkgs = append(met_pkgs, pkg.Metadata)
	}

	return met_pkgs
}
