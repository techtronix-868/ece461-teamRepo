package datastore

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mabaums/ece461-web/backend/models"
)

type InMemoryDatstore struct {
	packages    []models.Package
	initialized bool
}

func (md *InMemoryDatstore) initIfEmpty() {
	names := []string{"Foo package", "Bar Package", "Test Package", "Mark Package"}
	if !md.initialized {
		log.Printf("Initializing InMemoryDatastore")
		md.initialized = true
		for i := 0; i <= 40; i++ {
			pkg_meta := models.PackageMetadata{
				Name:    names[i%4],
				Version: fmt.Sprintf("1.0.%v", i),
				ID:      strconv.Itoa(i),
			}
			pkg_data := models.PackageData{
				Content:   "Content",
				JSProgram: "string",
			}
			pkg := models.Package{Data: pkg_data, Metadata: pkg_meta}
			md.packages = append(md.packages, pkg)
		}
	}
}

func (md *InMemoryDatstore) GetPackage(id string) (*models.Package, *models.Error) {
	md.initIfEmpty()
	log.Printf("Getting Package id: %v\n", id)
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

func packageMatch(pkg_met models.PackageMetadata, name string, version string) bool {
	if name != pkg_met.Name && name != "*" {
		return false
	}

	if version != pkg_met.Version && version != "" {
		return false
	}

	return true

}

func (md *InMemoryDatstore) ListPackages(offset int, pagesize int, name string, version string) []models.PackageMetadata {
	md.initIfEmpty()
	log.Printf("Listing packages offset: %v, pagesize: %v, name: %v, version: %v\n", offset, pagesize, name, version)
	found_packages := []models.PackageMetadata{}
	for _, pkg := range md.packages {
		if packageMatch(pkg.Metadata, name, version) {
			found_packages = append(found_packages, pkg.Metadata)
		}
	}
	return found_packages
}
