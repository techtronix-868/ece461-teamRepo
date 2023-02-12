package npm

import (
	nd "app/output"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Connect_npm struct {
	Package 	 string
	Version      string 
	Maintainers  int    
	Contributors int    
	License      string 
	Dependencies int    
	DevDeps      int    
	Releases     int    
	TestScript   bool   
	Commits      int    
	Downloads    int    
}


type Package struct {
	AnalyzedAt string `json:"analyzedAt"`
	Collected  struct {
		Metadata struct {
			Name        string `json:"name"`
			Scope       string `json:"scope"`
			Version     string `json:"version"`
			Description string `json:"description"`
			Keywords    []string `json:"keywords"`
			Date        string `json:"date"`
			Author      struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
			Publisher struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"publisher"`
			Maintainers []struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"maintainers"`
			Contributors []struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"contributors"`
			Repository struct {
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"repository"`
			Links struct {
				Npm       string `json:"npm"`
				Homepage  string `json:"homepage"`
				Repository string `json:"repository"`
				Bugs      string `json:"bugs"`
			} `json:"links"`
			License    string `json:"license"`
			Dependencies map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
			Releases []struct {
				From string `json:"from"`
				To   string `json:"to"`
				Count int    `json:"count"`
			} `json:"releases"`
			HasTestScript bool `json:"hasTestScript"`
		} `json:"metadata"`
		NPM struct {
			Downloads []struct {
				From string `json:"from"`
				To   string `json:"to"`
				Count int    `json:"count"`
			} `json:"downloads"`
			StarsCount int `json:"starsCount"`
		} `json:"npm"`
		Github struct {
			Homepage        string `json:"homepage"`
			StarsCount      int    `json:"starsCount"`
			ForksCount      int    `json:"forksCount"`
			SubscribersCount int    `json:"subscribersCount"`
			Issues          struct {
				Count      int `json:"count"`
				OpenCount  int `json:"openCount"`
				Distribution struct {
					OneHr    int `json:"3600"`
					ThreeHr  int `json:"10800"`
					NineHr   int `json:"32400"`
					OneDay   int `json:"97200"`
					ThreeDay int `json:"291600"`
					OneWk    int `json:"874800"`
					TwoWk    int `json:"2624400"`
					OneMo    int `json:"7873200"`
					ThreeMo  int `json:"23619600"`
					SixMo    int `json:"70858800"`
					OneYr    int `json:"212576400"`
				} `json:"distribution"`
				IsDisabled bool `json:"isDisabled"`
			} `json:"issues"`
			Contributors []struct {
				Username     string `json:"username"`
				CommitsCount int    `json:"commitsCount"`
			} `json:"contributors"`

		}`json:"github"`
	} `json:"collected"`
}

func (cn Connect_npm) Data(packageName string) *nd.NdJson {

	// The following makes an API call to NPM site and recieves JSON response.
	res1 := strings.Split(packageName,"/")
	packageName = res1[len(res1)-1]
	resp, err := http.Get(fmt.Sprintf("https://api.npms.io/v2/package/%s", packageName))
	if err != nil {
		
	}
	defer resp.Body.Close()

	// Marshallig JSON response onto the required struct
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err) 
	}

	var pkg Package

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal([]byte(body), &pkg)

	cn.Package = pkg.Collected.Metadata.Name
	cn.Version = pkg.Collected.Metadata.Version
	cn.Maintainers = len(pkg.Collected.Metadata.Maintainers)
	if (len(pkg.Collected.Metadata.Contributors) > 0){
		cn.Contributors = len(pkg.Collected.Metadata.Contributors)
	} else {
		cn.Contributors = len(pkg.Collected.Github.Contributors)
	}
	cn.License = pkg.Collected.Metadata.License
	cn.Dependencies = len(pkg.Collected.Metadata.Dependencies)
	cn.DevDeps = len(pkg.Collected.Metadata.DevDependencies)
	cn.Releases = len(pkg.Collected.Metadata.Releases)
	cn.TestScript = pkg.Collected.Metadata.HasTestScript
	cn.Commits = 0
	for _,s := range pkg.Collected.Github.Contributors{
		cn.Commits += s.CommitsCount
	}
	cn.Downloads = 0
	for _,s := range pkg.Collected.NPM.Downloads{
		cn.Downloads += s.Count
	}


	return cn.Score()
}

func (cn Connect_npm) Score() *nd.NdJson {
	

	overallScore:= 0.4*cn.get_responsivnesss()+0.1*cn.get_bus_factor() + 0.2*cn.get_License_score() + 0.1*cn.get_rampup_score() + 0.2 * cn.get_correctness()
	nd := new(nd.NdJson)
	nd=nd.DataToNd(cn.Package,overallScore,cn.get_rampup_score(),cn.get_bus_factor(),cn.get_responsivnesss(),cn.get_correctness(),cn.get_License_score())
	return nd
}

func Contains(sl []string, name string) bool {
	for _, value := range sl {
	   if value == name {
		  return true
	   }
	}
	return false
}

func (cn Connect_npm) get_License_score() float64 {
	cmpLicenses := []string{"Public Domain","MIT","X11","BSD-new","Apache 2.0","LGPLv2.1","LGPLv2.1+", "LGPLv3", "LGPLv3+"}

	if Contains(cmpLicenses,cn.License){
		return 1.0
	}
	
	return 0.0

}

func (cn Connect_npm) get_rampup_score() float64 {
	return float64(cn.DevDeps) / float64(cn.Dependencies)
}

func (cn Connect_npm) get_bus_factor() float64 {
	return float64(cn.Maintainers) / float64(cn.Contributors)
}

func (cn Connect_npm) get_correctness() float64{
	str1 := cn.Version
	res1 := strings.Split(str1, ".")
	major,_ := strconv.ParseFloat(res1[0],64)
	minor,_ := strconv.ParseFloat(res1[1],64)
	patch,_ := strconv.ParseFloat(res1[2],64)

	denominator := 1.0
	if major != 0 {
		denominator = major
	}
	if minor != 0 {
		denominator *= minor
	}
	if patch != 0 {
		denominator *= patch
	}
	
	score := 1 - ((major + minor + patch) / (denominator + 1))
	return float64(score)
}

func (cn Connect_npm) get_responsivnesss() float64{
	return float64(cn.Releases) / float64(cn.Commits)
}