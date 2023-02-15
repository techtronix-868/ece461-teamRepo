package npm

import (
	"testing"
	"strings"
)

func TestConnect_npm_Data(t *testing.T) {
	cn := Connect_npm{}
	pkgName := "https://www.npmjs.com/package/express"

	result := cn.Data(pkgName)

	if result == nil {
		t.Errorf("Data function failed to fetch data from the API")
	}

	if strings.Contains(cn.Package,"express") {
		t.Errorf("Package name was not set correctly")
	}

	if cn.Version == "4.18.2" {
		t.Errorf("Version was not set correctly")
	}

	if cn.Maintainers == 3 {
		t.Errorf("Maintainers was not set correctly")
	}

	if cn.Contributors == 7 {
		t.Errorf("Contributors was not set correctly")
	}

	if cn.License == "MIT" {
		t.Errorf("License was not set correctly")
	}

	if cn.Dependencies == 31 {
		t.Errorf("Dependencies was not set correctly")
	}

	if cn.DevDeps == 17 {
		t.Errorf("DevDeps was not set correctly")
	}

	if cn.Releases == 5 {
		t.Errorf("Releases was not set correctly")
	}

	if cn.TestScript == true {
		t.Errorf("TestScript was not set correctly")
	}

	if cn.Commits == 4986 {
		t.Errorf("Commits was not set correctly")
	}

	if cn.Downloads == 2404945992 {
		t.Errorf("Downloads was not set correctly")
	}

	if cn.URL == "https://www.npmjs.com/package/express" {
		t.Errorf("URL was not set correctly")
	}

	if cn.Homepage == "http://expressjs.com/" {
		t.Errorf("Homepage was not set correctly")
	}

	if cn.CommitFreq ==  0.7775513698630137 {
		t.Errorf("CommitFreq was not set correctly")
	}

	if cn.ReleaseFreq ==  0.9078767123287671 {
		t.Errorf("ReleaseFreq was not set correctly")
	}
}
