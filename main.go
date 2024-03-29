package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Base_url string
	Token    string
	Project  map[string]Project
}

type Project struct {
	Id      uint
	Path    string
	Depends []string `toml:"depends_on"`
	Ref     string
}

type BranchResponse struct {
	WebUrl string `json:"web_url"`
}

const PRIVATE_TOKEN_HEADER_NAME = "PRIVATE-TOKEN"

func main() {
	showHelp := flag.Bool("h", true, "Show help documentation")
	if *showHelp {
		show_help_docs()
		os.Exit(0)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unexpected HOME variable: We where unable to determine your HOME directory.")
		show_help_docs()
		os.Exit(1)
	}
	path := filepath.Join(home, "lazy-deploy.toml")
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Not found error: Config was not found.")
		show_help_docs()
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		show_help_docs()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Invalid argument count: Not enough arguments have been passed.")
		show_help_docs()
		os.Exit(1)
	}

	projectName := os.Args[1]
	branch := os.Args[2]

	var conf Config
	toml.Unmarshal(content, &conf)
	selectedProject, exists := conf.Project[projectName]
	if exists {
		create_branh(branch, conf, selectedProject)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stdout, "Project: \"%s\" was not found among your settings.\n", projectName)
	fmt.Fprintf(os.Stdout, "Available options are:\n")
	for k := range conf.Project {
		fmt.Fprintf(os.Stdout, "- %s\n", k)
	}
}

func show_help_docs() {
	fmt.Println(`LAZY DEPLOY:

	A tool to help you automate deploying a set of gitlab branches when you are too lazy.

USAGE:

	Simply invoke Lazy Deploy by specifying the project you want to deploy and a branch name.

		lazy-deploy <flags> <Project-Name> <New-Branch-Name>
	
	Available flags:

		-h	Display help documentation

CONFIG FILE:

	Lazy Deploy will read a configuration file written in named "lazy-deploy.toml" for keeping a reccord of each project and some
	other settings like your authentication token and your Gitlab instance's URL. Let's break down the config's content.

	There are 2 top level settings in your configuration.
		- base_url: The base URL of your gitlab instance. For example: "https://gitlab.com"
		- token: Your personal access token, please refer to the gitlab documentation: https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html

		Warning: Create a person access token in your "User Settings", this allows Lazy Deploy to read all your repositories, if you provide a project
		access token, then Lazy Deploy will not be able to access other projects.

	Now, we will define our project network. Here you will define each project you want to manage and which other projects does it depend on to be deployed.
		- id: Your project's ID, you will find this settings under Settings > General > Project ID.
		- path: If you are not able to retrieve your project's ID, you can provide a namespaced path to the project.

		  Consider the following example:

		      - A repository whose URL is: https://gitlab.com/test-user/test_project

		  A plausible configuration for this repository could be.

			[Project.Test]
			# Id = 123
			Path = "test-user/test_project"
			ref = "main"
			depends_on = []

		  In this case I did not provide an Id, you can include it, but *ID will take precedence over Path when building your request for creating
		  the new branch*. Also, the Gitlab API recommends that this path is URL encoded, we take care of this, so you can declare the path
		  as you normally would.

		- ref: The "base" branch which you project will use for creating new branches.
		- depends_on: This is an array of projects, this means that when "X" project is deployed, all the projects it depends on will also be deployed.

	CONFIG FILE LOCATION:

		You must place your configuration file in your $HOME directory in UNIX operating systems, or in %\USERPROFILE% in Windows.
		This file must be named: "lazy-deploy.toml"

	CONFIG FILE EXAMPLE:

		Here's an example configuration file.

			base_url = "https://gitlab.com"
			token = "glpat-123"

			[Project.ServiceA]
			id = 123
			ref = "main"
			depends_on = ["ServiceB"]

			[Project.ServiceB]
			id = 321
			ref = "main"
			depends_on = []

			[Project.ServiceC]
			Path = "test-username/test_project"
			ref = "main"
			depends_on = []
		
	SEE ALSO:
		- Tom's Obvious Minimal Language's Official Page: https://toml.io/en/

	NOTES:
	`)
}

func build_new_branch_endpoint(conf Config, project Project, branch string) string {
	if project.Id == 0 {
		urlEncodedPath := url.PathEscape(project.Path)
		return fmt.Sprintf("projects/%s/repository/branches?branch=%s&ref=%s", urlEncodedPath, branch, project.Ref)
	}
	return fmt.Sprintf("projects/%d/repository/branches?branch=%s&ref=%s", project.Id, branch, project.Ref)
}

func create_branh(branch string, conf Config, project Project) {
	apiEndpoint := fmt.Sprintf("%s/api/v4", conf.Base_url)
	newBranchEndpoint := build_new_branch_endpoint(conf, project, branch)
	requestEndpoint := fmt.Sprintf("%s/%s", apiEndpoint, newBranchEndpoint)

	newBranchRequest, _ := http.NewRequest("POST", requestEndpoint, nil)
	newBranchRequest.Header.Add(PRIVATE_TOKEN_HEADER_NAME, conf.Token)
	response, _ := http.DefaultClient.Do(newBranchRequest)

	if response.StatusCode == 201 {
		defer response.Body.Close()
		decoder := json.NewDecoder(response.Body)
		var branch BranchResponse
		decoder.Decode(&branch)
		fmt.Printf("%s\n", branch.WebUrl)
	} else {
		fmt.Fprintf(os.Stderr, "An error has ocurred while creating your branch: %s\n", response.Status)
		os.Exit(1)
	}

	if len(project.Depends) > 0 {
		for _, v := range project.Depends {
			create_branh(branch, conf, conf.Project[v])
		}
	}
}
