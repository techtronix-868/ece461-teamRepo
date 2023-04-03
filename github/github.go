package github

import (
    "context"
    "fmt"
	
	"strings"
	"strconv"
    "github.com/machinebox/graphql"

	"math"
	"os"
	nd "app/output"
	lg "app/lg"
	git "app/git"

)


var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")


//Fuction to round off numbers
func roundFloat(val float64, prec uint) float64 {
    ratio := math.Pow(10, float64(prec))
    return math.Round(val*ratio) / ratio
}


//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

//Function to get total number of commits in a repository
func Get_com(owner string, name string) int {
    graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	lg.Init(os.Getenv("LOG_FILE"))
	
    graphqlRequest := graphql.NewRequest(`
		query Get_commits($own: String!, $repo: String!) {
			repository(owner:$own, name:$repo) {
				defaultBranchRef {
					target {
						... on Commit {
							history {
								totalCount
							}
						}
					}
				}
			}
			
		}
    `)

	// repository(owner:$own, name:$repo) {
	// 	object(expression:"master") {
	// 		... on Commit {
	// 		history {
	// 			totalCount
	// 		}
	// 		}
	// 	}
	// 	}

	graphqlRequest.Var("own",owner)
	graphqlRequest.Var("repo",name)



	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        lg.ErrorLogger.Println("Unable to get commits through GrpahQL API in github.go")
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	strs := strings.SplitAfter(str, "totalCount:")
	strss := fmt.Sprint(strs[1])
	strsss := strings.Split(strss,"]")

	commits,err := strconv.Atoi(strsss[0])

	if err != nil {
		lg.ErrorLogger.Println("Unable to get round number of commits")
		return 0
	}	
	
	lg.InfoLogger.Println("Setting number of commits : ",commits)
	return commits

}

func Get_releases(owner string, name string) int {

	graphqlClient := graphql.NewClient("https://api.github.com/graphql")

	lg.Init(os.Getenv("LOG_FILE"))
	
	
    graphqlRequest := graphql.NewRequest(`
	query Get_commits($own: String!, $repo: String!){
		repository(name:$repo, owner: $own) {
			releases {
			totalCount
			}
		}
		}
		
    `)

	graphqlRequest.Var("own",owner)
 	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        lg.ErrorLogger.Println("Unable to reach releases through GrpahQL API in github.go")
		return 0
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "totalCount:")
	//fmt.Println(strs[1])
    parse_2 := fmt.Sprint(parse_1[1])
    parse_3 := strings.Split(parse_2,"]")

	rels,err := strconv.Atoi(parse_3[0])

	if err != nil {
		lg.ErrorLogger.Println("Unable to get round number of releases")
		os.Exit(1)
	}

	lg.InfoLogger.Println("Setting number of releases : ",rels)

	return rels	
		
}


func ScoreResponsiveness(owner string,repo string) float64 {

	lg.Init(os.Getenv("LOG_FILE"))

	com := Get_com(owner,repo)
	releases := Get_releases(owner,repo)

	score := float64(releases) / float64((com + 1))
	if score > 1{
        score = 1
	} else if score < 0{
		score = 0
	}

	lg.InfoLogger.Println("Finding responsiveness score : ",score)

	return roundFloat(score,2)
        


}

func Get_assignees(owner string, name string) int{

	lg.Init(os.Getenv("LOG_FILE"))
	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    graphqlRequest := graphql.NewRequest(`
	query Get_commits($own: String!, $repo: String!){
		repository(name:$repo, owner: $own) {
			assignableUsers {
			totalCount
			}
		}
		}
		
    `)

	graphqlRequest.Var("own",owner)
 	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        lg.ErrorLogger.Println("Unable to reach assignees through GrpahQL API in github.go")
		return 0
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "totalCount:")
	//fmt.Println(strs[1])
    parse_2 := fmt.Sprint(parse_1[1])
    parse_3 := strings.Split(parse_2,"]")

	assign,err := strconv.Atoi(parse_3[0])

	if err != nil {
		lg.ErrorLogger.Println("Unable to get round number of assignees")
		os.Exit(1)
	}

	lg.InfoLogger.Println("Setting number of assignees : ",assign)

	return assign


}

func Get_contributors(owner string, name string) int{

	lg.Init(os.Getenv("LOG_FILE"))

	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    graphqlRequest := graphql.NewRequest(`

	query Get_commits($own: String!, $repo: String!){
		repository(owner:$own, name: $repo) {
			id
			name
			mentionableUsers {
			totalCount
			}
		}
		}
		
    `)

	graphqlRequest.Var("own",owner)
 	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
        lg.ErrorLogger.Println("Unable to reach assignees through GrpahQL API in github.go")
		os.Exit(1)
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "totalCount:")
	//fmt.Println(strs[1])
    parse_2 := fmt.Sprint(parse_1[1])
    parse_3 := strings.Split(parse_2,"]")

	contributors,err := strconv.Atoi(parse_3[0])

	if err != nil {
		lg.ErrorLogger.Println("Unable to get round number of contributors")
		os.Exit(1)
	}

	lg.InfoLogger.Println("Setting number of contributors : ",contributors)

	return contributors


}

func ScoreBusFactor(owner string, repo string) float64 {

	lg.Init(os.Getenv("LOG_FILE"))

	assign := Get_assignees(owner,repo)
	contributors := Get_contributors("nullivex","nodist")

	score := float64(assign) / float64((contributors + 1))

	if score > 1{
        score = 1
	} else if score < 0{
		score = 0
	}

	lg.InfoLogger.Println("Finding responsiveness score : ",score)
	return roundFloat(score,2)
        


}


func get_dependancies(owner string, name string) int{

	lg.Init(os.Getenv("LOG_FILE"))

	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    graphqlRequest := graphql.NewRequest(`
	query Get_commits($own: String!, $repo: String!){
	repository(owner:$own, name:$repo) {
		dependencyGraphManifests {
		edges {
			node {
				dependenciesCount
			}
		}
		}
	  }
	}
		
    `)

	graphqlRequest.Var("own",owner)
 	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
	graphqlRequest.Header.Set("Accept", "application/vnd.github.hawkgirl-preview+json")
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		lg.ErrorLogger.Println("Unable to get depandancies through GrpahQL API in github.go")
		os.Exit(1)
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "dependenciesCount:")
	//fmt.Println(parse_1[1])
    parse_2 := fmt.Sprint(parse_1[1])
    parse_3 := strings.Split(parse_2,"]")

	dependencies,err := strconv.Atoi(parse_3[0])

	if err != nil {
		lg.ErrorLogger.Println("Unable to get round number of depandancies")
		os.Exit(1)
	}

	lg.InfoLogger.Println("Setting number of depandencies : ",dependencies)

	return dependencies


}


func get_devDep(owner string, name string) int {

	lg.Init(os.Getenv("LOG_FILE"))

	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    graphqlRequest := graphql.NewRequest(`
	query Get_commits($own: String!, $repo: String!){
		repository(owner:$own, name:$repo) {
		  dependencyGraphManifests {
			edges {
			  node {
				dependencies {
				  totalCount
				  nodes {
					hasDependencies
				  }
				}
			  }
			}
		  }
		}
	  }
		
    `)

	graphqlRequest.Var("own",owner)
	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
	graphqlRequest.Header.Set("Accept", "application/vnd.github.hawkgirl-preview+json")
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		lg.ErrorLogger.Println("Unable to get Dev depandancies through GrpahQL API in github.go")
        os.Exit(1)
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "totalCount:")
	//fmt.Println(parse_1[0])

    parse_2 := fmt.Sprint(parse_1[0])
	devDep := strings.Count(parse_2, "true")

	lg.InfoLogger.Println("Setting number of Dev depandancies : ",devDep)

    
	if devDep > 0 {
		return devDep
	}else {
		return 0
	}
	

}

func scoreRampUp(owner string,repo string) float64 {

	lg.Init(os.Getenv("LOG_FILE"))

	dependencies := get_dependancies("nullivex","nodist")
	devDep := get_devDep(owner,repo)

	score := float64(devDep) / float64((dependencies + 1))

	if score > 1{
        score = 1
	} else if score < 0{
		score = 0

	}

	lg.InfoLogger.Println("Finding RampUp score : ",score)

	return roundFloat(score,2)
        


}
func get_License(owner string, name string) string{

	lg.Init(os.Getenv("LOG_FILE"))

	url := fmt.Sprintf("https://github.com/%s/%s",owner,name)

	defer func() {
		if err := recover(); err != nil {
		
		}
	}()

	if git.Clone(url){
		return "present"
	}


	

	// graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    // graphqlRequest := graphql.NewRequest(`
	// query Get_commits($own: String!, $repo: String!){
	// 		repository(name: $repo, owner: $own) {
	// 		  licenseInfo {
	// 			name
	// 		  }
			  
	// 	}
	// }
		
    // `)

	// graphqlRequest.Var("own",owner)
 	// graphqlRequest.Var("repo",name)


	// graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    // var graphqlResponse interface{}
    // if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
	// 	lg.ErrorLogger.Println("Unable to get License through GrpahQL API in github.go")
    //     os.Exit(1)
    // }

	// str := fmt.Sprint(graphqlResponse)
	// fmt.Println(str)

	// found_license := strings.Count(str, "name:")
	// fmt.Println(found_license)
	// if found_license == 0{
	// 	lg.WarningLogger.Println("No license Found")
	// 	return "No license found"
	// }

	// parse_1 := strings.SplitAfter(str, "name:")
	// //fmt.Println(parse_1[0])
    // parse_2 := fmt.Sprint(parse_1[1])
    // parse_3 := strings.Split(parse_2,"]")
	// parse_4 := fmt.Sprint(parse_3[0])
	// parse_5 := strings.Split(parse_4,"License")
	// parse_6 := fmt.Sprint(parse_5[0])
	// parse_7 := strings.Split(parse_6," ")

	// //fmt.Println(parse_5[0])

	//lg.InfoLogger.Println("License is  : ",parse_7[0])


	return "not_present"




}


func scoreLicense(owner string, repo string) float64{
	lg.Init(os.Getenv("LOG_FILE"))

	license := get_License(owner,repo)
	
	if "present" == license {
		  lg.InfoLogger.Println("LicenseScore  is  : ",1)
		  return 1.0
	}
	lg.InfoLogger.Println("LicenseScore is  : ",0)
	return 0.0


}

func get_tag(owner string, name string) string {

	lg.Init(os.Getenv("LOG_FILE"))

	graphqlClient := graphql.NewClient("https://api.github.com/graphql")
	
	
    graphqlRequest := graphql.NewRequest(`
	query Get_commits($own: String!, $repo: String!){
			repository(name: $repo, owner: $own) {
			  releases(last: 1) {
				nodes{
					tagName
				}
			  }
			  
		}
	}
		
    `)

	graphqlRequest.Var("own",owner)
 	graphqlRequest.Var("repo",name)


	graphqlRequest.Header.Set("Authorization", "Bearer " + GITHUB_TOKEN)
    var graphqlResponse interface{}
    if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		lg.ErrorLogger.Println("Unable to get Repository version through GrpahQL API in github.go")
        os.Exit(1)
    }

	str := fmt.Sprint(graphqlResponse)
	//fmt.Println(str)

	parse_1 := strings.SplitAfter(str, "tagName:")
	if len(parse_1) == 1{
		return "No version"
	}
	//fmt.Println(parse_1[1])
    parse_2 := fmt.Sprint(parse_1[1])
    parse_3 := strings.Split(parse_2,"]")
	//fmt.Println(parse_3[0])
	parse_4 := fmt.Sprint(parse_3[0])
	// parse_5 := strings.Split(parse_4,"License")
	// parse_6 := fmt.Sprint(parse_5[0])
	parse_7 := strings.Split(parse_4," ")

	lg.InfoLogger.Println("Setting the repository version : ",parse_7[0])


	return parse_7[0]
	//fmt.Println(parse_7[0])


}


func scoreCorrectness(owner string, repo string) float64{

	lg.Init(os.Getenv("LOG_FILE"))

	version := get_tag(owner,repo)

	if version == "No version"{
		return 0.0
	}

	res1 := strings.Split(version, ".")
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

	lg.InfoLogger.Println("Finding Correctness score : ",float64(score))
	return float64(score)
}

func Score(URL string) *nd.NdJson {

	lg.Init(os.Getenv("LOG_FILE"))
	
	cuttingByTwo := strings.FieldsFunc(URL, func(r rune) bool {
		if r == '/' {
			return true
		}
		return false
	})

	//fmt.Println(cuttingByTwo)

	var owner string = cuttingByTwo[2]
	var repo string = cuttingByTwo[3]


	overallScore := 0.4*ScoreResponsiveness(owner,repo) + 0.1*ScoreBusFactor(owner,repo) + 0.2*scoreLicense(owner, repo) + 0.1*scoreRampUp(owner,repo) + 0.2 * scoreCorrectness(owner,repo)
	lg.InfoLogger.Println("Finding overall score : ",overallScore)
	nd := new(nd.NdJson)
	nd=nd.DataToNd(URL,overallScore,scoreRampUp(owner,repo),ScoreBusFactor(owner,repo),ScoreResponsiveness(owner,repo),scoreCorrectness(owner,repo),scoreLicense(owner,repo))

	return nd
}

