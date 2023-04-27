package packager

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mabaums/ece461-web/backend/models"
	log "github.com/sirupsen/logrus"
)

type PackageJSON struct {
	Name    string
	Version string
}

func GetPackageJson(url string) (*models.PackageMetadata, error) {
	tempDir, err := os.MkdirTemp(".", "*")

	if err != nil {
		log.Errorf("Error creating temporary folder %v", err)
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	err = Clone(tempDir, url)

	if err != nil {
		return nil, err
	}

	return ReadPackageJson(tempDir)
}

func Clone(dir string, url string) error {
	log.Infof("Cloning %v into %v", url, dir)
	cmd := exec.Command("git", "clone", url, dir)
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error Cloning: %v in Dir: %v Err: %v", url, dir, err) // Maybe no need to be Fatal?
	}
	return err
}

func ReadPackageJson(dir string) (*models.PackageMetadata, error) {
	log.Infof("Reading package.json in %v", dir)
	var metadata models.PackageMetadata

	content, err := ioutil.ReadFile(filepath.Join(dir, "package.json"))

	if err != nil {
		log.Errorf("No package.json found in %v", dir)
		return nil, err
	}

	var pkgJSON PackageJSON

	err = json.Unmarshal(content, &pkgJSON)

	if err != nil {
		log.Errorf("package.json is invalid: %v", err)
		return nil, err
	}
	log.Infof("Parsed package.json %+v", pkgJSON)

	metadata = models.PackageMetadata{
		Name:    pkgJSON.Name,
		Version: pkgJSON.Version,
	}

	return &metadata, nil

}
