package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Data represents the entries in metadata.yaml
type Data struct {
	Files []struct {
		Path        string   `yaml:"path"`
		Description string   `yaml:"description"`
		Services    []string `yaml:"services"`
		Operations  []string `yaml:"operations"`
	} `yaml:"files"`
}

// Global struct to hold metadata for gov2 and gov2/<service> folders
var files Data

func populateFiles(debug bool, path string) error {
	filePrefix := "https://raw.githubusercontent.com/awsdocs/aws-doc-sdk-examples/master/"
	debugPrint(debug, "Getting metadata from "+filePrefix+path)
	results, err := http.Get(filePrefix + path)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(results.Body)
	text := buf.String()

	err = yaml.Unmarshal([]byte(text), &files)
	if err != nil {
		fmt.Println("Got an error unmarshalling " + filePrefix + path)
		return err
	}

	return nil
}

// RepoTree represents the files in a repo
// See https://docs.github.com/en/free-pro-team@latest/rest/reference/git#trees
type RepoTree struct {
	Sha  string `json:"sha"` // Like a checksum for the request
	URL  string `json:"url"` // The API URL for the request
	Tree []struct {
		Path string `json:"path"` // The path to the directory/file after https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/
		Mode string `json:"mode"` // The file permissions. Most folders are ; files 100644
		Type string `json:"type"` // "blob" for files; "tree" for directories
		Sha  string `json:"sha"`  // Like a checksum for the directory/file
		Size int    `json:"size"` // The size, in bytes, of files
		URL  string `json:"url"`  // The API URL for the file
	} `json:"tree"`
	Truncated bool `json:"truncated"` // Whether the request was truncated (more to come)
}

func debugPrint(debug bool, s string) {
	if debug {
		fmt.Println(s)
	}
}

func isValueLanguage(debug bool, lang string) bool {
	for _, l := range globalConfig.Languages {
		if lang == l {
			return true
		}
	}

	return false
}

// This only works for Go currently
func getNameForSvcDir(debug bool, dir string) string {
	switch dir {
	case "cloudwatch":
		return "Amazon CloudWatch"
	case "dynamodb":
		return "Amazon DynamoDB"
	case "ec2":
		return "Amazon EC2"
	case "iam":
		return "IAM"
	case "kms":
		return "AWS KMS"
	case "s3":
		return "Amazon S3"
	case "sns":
		return "Amazon SNS"
	case "sqs":
		return "Amazon SQS"
	case "ssm":
		return "AWS Systems Manager"
	case "sts":
		return "AWS STS"
	default:
		return ""
	}
}

func createSvcDir(debug bool, language string, dir string, outDir string, translation string) error {
	name := getNameForSvcDir(debug, dir)
	if name == "" {
		msg := dir + " is not recognized as a valid service directory"
		return errors.New(msg)
	}

	outputDir := outDir + "/" + dir
	outputFile := outputDir + "/"

	// Create directory for service
	debugPrint(debug, "Creating service directory "+outputDir)
	err := os.Mkdir(outDir+"/"+dir, 0755)
	if err != nil {
		return err
	}

	switch translation {
	case "index":
		// Create _index.md files for top-level service topics
		outputFile += "_index.md"
		break
	default:
		// Otherwise, create <service>_index.md
		outputFile += dir + "_index.md"
	}

	// Create topic file in that directory
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	content := ""

	switch translation {
	case "index":
		// Here's what we write to [<service>]_index.md
		content = "---\n" +
			"title: \"" + name + " Examples\"\n" +
			"linkTitle: \"" + name + "\"\n" +
			"weight: 3\n" +
			"---\n" +
			"\n" +
			"This section contains code examples for " + name + " using version 2 of the AWS SDK for Go.\n"

		break
	default:
		content = "This section contains code examples for " + name + " using version 2 of the AWS SDK for Go.\n"
		// Add after updating metadata.yaml with description, operations:
		//   1. Get all paths that don't end in _test.go -> ## path
		//   2. Get description for that entry           -> Description
		break
	}

	_, err = w.WriteString(content)
	if err != nil {
		fmt.Println("Got an error writing to " + outputFile)
		return err
	}

	/* If we ever want to add the contents of the README.md file in that folder:
	filePrefix := "https://raw.githubusercontent.com/awsdocs/aws-doc-sdk-examples/master/"
	path := filePrefix + language + "/" + dir + "/" + "README.md"
	results, err := http.Get(path)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(results.Body)
	text := buf.String()

	_, err = w.WriteString(text)
	if err != nil {
		fmt.Println("Got an error writing README.md to " + outputFile)
		return err
	}
	*/

	w.Flush()

	return nil
}

func createOperationDir(debug bool, dir string, subdir string, outDir string) error {
	outputDir := outDir + "/" + dir + "/" + subdir

	// Create directory
	debugPrint(debug, "Creating operation directory "+outputDir)
	err := os.Mkdir(outputDir, 0755)

	return err
}

func getNameFromMetadata(debug bool, subdir string, file string) (string, error) {
	// Wade through files
	// If subdir/file == path, return operation[0], if it exists
	for _, f := range files.Files {
		debugPrint(debug, "Looking at metadata path: "+f.Path)

		// Until we get operations for all Go v2 service metadata files
		if f.Path == subdir+"/"+file {
			if nil != f.Operations && f.Operations[0] != "" {
				return f.Operations[0], nil
			}
		}

		debugPrint(debug, "No operations for "+subdir+"/"+file)
		// Just return the filename, with .go replace by .md
		file = strings.Replace(file, ".go", ".md", 1)
		return file, nil
	}

	return "", nil
}

func createFile(debug bool, dir string, subdir string, fileName string, path string, outDir string, translation string) error {
	filePrefix := "https://raw.githubusercontent.com/awsdocs/aws-doc-sdk-examples/master/"
	linkPrefix := "https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/gov2/"
	linkFile := linkPrefix + dir + "/" + subdir + "/" + fileName

	filePath := filePrefix + path

	// Is there a README.md file in the same directory?
	//   filePath is something like:
	//     https://raw.githubusercontent.com/awsdocs/aws-doc-sdk-examples/master/gov2/cloudwatch/DescribeAlarms/DescribeAlarmsv2.go
	// So split the string by '/', get the last part, and replace it with README.md.
	mdFilePath := strings.Replace(filePath, fileName, "README.md", 1)

	// Snarf README.md and use it to start building the local file
	results, err := http.Get(mdFilePath)
	if err != nil {
		return err
	}

	outFileName := outDir + "/" + dir + "/" + subdir

	// The file to create
	switch translation {
	case "index":
		outFileName += "/_index.md"
		break
	case "metadata":
		name, err := getNameFromMetadata(debug, subdir, fileName)
		if err != nil {
			return err
		}

		debugPrint(debug, "Name from metadata: "+name)

		outFileName += "/" + name
		break
	default: // none
		fileName = strings.Replace(fileName, ".go", ".md", 1)
		outFileName += "/" + fileName
	}

	debugPrint(debug, "Creating output file: "+outFileName)

	// Read each line of readme file from repo
	scanner := bufio.NewScanner(results.Body)

	f, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Got an error opening " + outFileName)
		return err
	}

	_, err = f.Write([]byte(fmt.Sprintf(`---
title: "%s"
---
`, strings.TrimSuffix(fileName, ".go"))))
	if err != nil {
		fmt.Println("Got an error writing to " + outFileName)
		return err
	}

	defer f.Close()

	// Dump contents of readme into new MD file
	for scanner.Scan() {
		// Convert third-level heading to top-level heading
		text := strings.Replace(scanner.Text(), "##", "#", 1)
		_, err = f.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Got an error writing to " + outFileName)
			return err
		}
	}

	// Add heading for code
	_, err = f.Write([]byte("\n## Source code\n\n"))
	if err != nil {
		fmt.Println("Got an error writing to " + outFileName)
		return err
	}

	// Get contents of repo file
	resp, err := http.Get(filePath)
	if err != nil {
		return err
	}

	scanner = bufio.NewScanner(resp.Body)

	// Start of code
	_, err = f.Write([]byte("```go\n"))
	if err != nil {
		fmt.Println("Got an error writing to " + outFileName)
		return err
	}

	for scanner.Scan() {
		// Strip out any snippet tags
		if !strings.Contains(scanner.Text(), "snippet-") {
			_, err = f.Write([]byte(scanner.Text() + "\n"))
			if err != nil {
				fmt.Println("Got an error writing to " + outFileName)
				return err
			}
		}
	}

	coda := "```\n" +
		"\n" +
		"See the [complete example in GitHub](" + linkFile + ")."

	// End of code
	_, err = f.Write([]byte(coda))
	if err != nil {
		fmt.Println("Got an error writing to " + outFileName)
		return err
	}

	return nil
}

func processFiles(debug bool, language string, input string, outDir string, translation string) error {
	// We're only looking at gov2 paths for now
	if language != "gov2" {
		return errors.New("Sorry, only gov2 is currently supported as a language")
	}

	// Unmarshal string into tree
	var repoTree RepoTree

	json.Unmarshal([]byte(input), &repoTree)

	svcDir := ""
	subDir := ""

	for _, leaf := range repoTree.Tree {
		// leaf.Path for a file looks like:
		//   gov2/cloudwatch/CreateCustomMetric/CreateCustomMetricv2.go

		if leaf.Path == language {
			debugPrint(debug, "Processing "+language+" root path: "+leaf.Path)
			if translation == "metadata" {
				// Read metadata so we can name output files correctly
				debugPrint(debug, "Processing "+leaf.Path+"/metadata.yaml")
				err := populateFiles(debug, leaf.Path+"/metadata.yaml")
				if err != nil {
					return err
				}
			}
		}

		// Split up the path '/'
		parts := strings.Split(leaf.Path, "/")
		// parts[0] == gov2
		// parts[1] == cloudwatch
		// parts[2] == CreateCustomeMetric
		// parts[3] == CreateCustomMetricv2.go

		if parts[0] != language {
			continue
		}

		switch len(parts) {
		case 2:
			// We have something like:
			//     gov2/cloudwatch or gov2/.gitignore
			// so skip files
			if leaf.Type == "blob" {
				break
			}

			if svcDir != parts[1] {
				debugPrint(debug, "Found new service directory: "+parts[1])

				if translation == "metadata" {
					// Read metadata so we can:
					// * name operation topic files correctly when metadata-based naming scheme is specified
					// * create list of code examples for service-level topic
					debugPrint(debug, "Processing "+parts[0]+"/"+parts[1]+"/metadata.yaml")
					err := populateFiles(debug, parts[0]+"/"+parts[1]+"/metadata.yaml")
					if err != nil {
						return err
					}
				}

				err := createSvcDir(debug, language, parts[1], outDir, translation)
				if err != nil {
					return err
				}

				svcDir = parts[1]
			}

			break
		case 3:
			// Now we have entries like:
			//     gov2/cloudwatch/CreateCustomMetric
			// Make sure we trap any files at this level
			// Files are of type "blob"
			if leaf.Type == "blob" {
				break
			}

			if subDir != parts[2] {
				debugPrint(debug, "Found new operation directory: "+parts[2])

				err := createOperationDir(debug, svcDir, parts[2], outDir)
				if err != nil {
					return err
				}

				subDir = parts[2]
			}

			break
		case 4:
			// Now we have entries like:
			//     gov2/cloudwatch/CreateCustomMetric/CreateCustomMetricv2.go
			//       parts[0] == gov2
			//       parts[1] == cloudwatch
			//       parts[2] == CreateCustomeMetric
			//       parts[3] == CreateCustomMetricv2.go
			//     gov2/cloudwatch/CreateCustomMetric/CreateCustomMetricv2_test.go
			//     gov2/cloudwatch/CreateCustomMetric/config.json
			// We only care about *.go, but not _test.go
			if strings.Contains(parts[3], "_test.go") {
				break
			}

			if strings.Contains(parts[3], ".go") {
				// For paths like:
				//     gov2/cloudwatch/CreateCustomMetric/CreateCustomMetricv2.go
				err := createFile(debug, svcDir, subDir, parts[3], leaf.Path, outDir, translation)
				if err != nil {
					return err
				}
			}

			break

		default:
			break
		}
	}

	return nil
}

// Config defines the configuration values from config.json
type Config struct {
	UserName    string   `json:"UserName"`
	Languages   []string `json:"Languages"`
	Language    string   `json:"Language"`
	OutDir      string   `json:"OutDir"`
	Translation string   `json:"Translation"`
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration() error {
	content, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("    go run TraverseRepoFiles.go -u NAME -l Language -o OUTPUT-DIR [-t none | index (default) | metadata] [-d] [-h] [-f]")
	fmt.Println(" where:")
	fmt.Println("    NAME      is the name of the GitHub user used to the GitHub API")
	fmt.Println("              the default is the value of UserName in config.json")
	fmt.Println("    LANGUAGE  is the programming language directory in our GitHub repo, such as gov2")
	fmt.Println("              the default is the value of Language in config.json")
	fmt.Println("    OUT-DIR   specifies where the code example topics are saved")
	fmt.Println("              the default is the value of OutDir in config.json")
	fmt.Println("    -d        (debug) displays additional debugging information")
	fmt.Println("    -h        (help) displays this error message and quits")
	fmt.Println("    -f        (fake it) displays the response as JSON and quits")
	fmt.Println("")
}

func main() {
	err := populateConfiguration()
	if err != nil {
		globalConfig.UserName = ""
		globalConfig.Language = ""
		globalConfig.OutDir = ""
		globalConfig.Translation = "index"
	}

	userName := flag.String("u", globalConfig.UserName, "Your GitHub user name")
	language := flag.String("l", globalConfig.Language, "The language in the GitHub repo")
	outDir := flag.String("o", globalConfig.OutDir, "Root directory where the output files are created")
	translation := flag.String("t", globalConfig.Translation, "Whether to translate source filename -> destination filename; valid values are none (just use existing filename), index (default, translate everything to _index.md), or metadata (use name from metadata.yaml)")
	debug := flag.Bool("d", false, "Whether to barf out more info. False by default.")
	fakeIt := flag.Bool("f", false, "Whether to just barf out the response")
	help := flag.Bool("h", false, "Displays usage and quits")
	flag.Parse()

	if *help {
		usage()
		return
	}

	if *userName == "" || *language == "" || *outDir == "" {
		usage()
		return
	}

	if *debug {
		fmt.Println("Debugging enabled")
		fmt.Println("User:        " + *userName)
		fmt.Println("Language:    " + *language)
		fmt.Println("Output dir:  " + *outDir)
		fmt.Println("Translation: " + globalConfig.Translation)
	}

	// We overload language in the toolchain
	if *language == "go" {
		*language = "gov2"
	}

	if !isValueLanguage(*debug, *language) {
		fmt.Println("The language directory " + *language + " is not present in the repo")
		return
	}

	if !(*translation == "none" || *translation == "index" || *translation == "metadata") {
		fmt.Println("The translation value " + *translation + " is neither none, index, nor metadata")
		return
	}

	gitHubURL := "https://api.github.com"
	query := gitHubURL + "/repos/awsdocs/aws-doc-sdk-examples/git/trees/master?recursive=1"

	debugPrint(*debug, "Querying: ")
	debugPrint(*debug, query)

	jsonData := ""
	jsonValue, _ := json.Marshal(jsonData)

	request, err := http.NewRequest("GET", query, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("Got an error creating HTTP request:")
		fmt.Println(err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/vnd.github.v3+json")

	request.SetBasicAuth(*userName, "")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return
	}

	data, _ := ioutil.ReadAll(response.Body)

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, data, "", "\t")
	if error != nil {
		fmt.Println("Got an error indenting JSON bytes:")
		fmt.Println(err)
		return
	}

	// If tryit set, just barf out response
	if *fakeIt {
		fmt.Println(string(prettyJSON.Bytes()))
		return
	}

	err = processFiles(*debug, *language, string(prettyJSON.Bytes()), *outDir, *translation)
	if err != nil {
		fmt.Println("Got an error processing files:")
		fmt.Println(err)
	}
}
