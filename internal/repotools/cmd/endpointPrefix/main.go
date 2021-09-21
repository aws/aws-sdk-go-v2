package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	repotools "github.com/awslabs/aws-go-multi-module-repository-tools"
)

var (
	atOnce      int
	rootPath    string
	pathRelRoot bool
)

var (
	modelGlobPattern       string
	endpointPrefixFilename string
)

func init() {
	flag.StringVar(&modelGlobPattern, "m", "",
		"The `glob pattern` of the API models.")

	flag.StringVar(&endpointPrefixFilename, "o", "",
		"The `endpoints prefix` path to write output to.")
}

const stdEndpointPrefixPath = `codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoint-prefix.json`

func run() (err error) {
	flag.Parse()

	if len(modelGlobPattern) == 0 {
		return fmt.Errorf("model glob pattern not provided")
	}

	if len(endpointPrefixFilename) == 0 {
		repoRoot, err := repotools.FindRepoRoot("")
		if err != nil {
			return fmt.Errorf("failed to find root of repository, %v", err)
		}
		endpointPrefixFilename = filepath.Join(repoRoot, filepath.FromSlash(stdEndpointPrefixPath))
	}

	filepaths, err := filepath.Glob(modelGlobPattern)
	if err != nil {
		return fmt.Errorf("failed to enumerate models %v, %w", modelGlobPattern, err)
	}

	prefixes := make(map[string]string, len(filepaths))
	for _, p := range filepaths {
		model, err := loadAPIModel(p)
		if err != nil {
			return err
		}

		model.Metadata.ServiceID = getServiceID(model)
		if _, ok := exclueServiceID[model.Metadata.ServiceID]; ok {
			log.Println("Ignoring", model.Metadata.ServiceID)
			continue
		}

		if len(model.Metadata.EndpointPrefix) == 0 {
			return fmt.Errorf("endpoint prefix missing for %v", p)
		}

		prefixes[model.Metadata.ServiceID] = model.Metadata.EndpointPrefix
	}

	f, err := os.Create(endpointPrefixFilename)
	if err != nil {
		return fmt.Errorf("unable to create %v, %w", endpointPrefixFilename, err)
	}
	defer func() {
		fErr := f.Close()
		if err == nil && fErr != nil {
			err = fmt.Errorf("failed to close endpoints prefix file, %w", fErr)
		}
	}()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(prefixes); err != nil {
		return fmt.Errorf("failed to encode endpoint prefixes, %w", err)
	}

	return nil
}

func loadAPIModel(filename string) (model apiModel, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return apiModel{}, fmt.Errorf("failed to load model, %w", err)
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&model); err != nil {
		return apiModel{}, fmt.Errorf("failed to decode model %v, %w", filename, err)
	}

	return model, nil
}

type apiModel struct {
	Metadata struct {
		ServiceID           string
		ServiceFullName     string
		ServiceAbbreviation string
		EndpointPrefix      string
	} `json:"metadata"`
}

func main() {
	if err := run(); err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}
}

var exclueServiceID = map[string]struct{}{
	"SimpleDB":     {},
	"ImportExport": {},
}

var serviceIDRegex = regexp.MustCompile("[^a-zA-Z0-9 ]+")
var prefixDigitRegex = regexp.MustCompile("^[0-9]+")

func getServiceID(model apiModel) string {
	if len(model.Metadata.ServiceID) > 0 {
		return strings.Title(model.Metadata.ServiceID)
	}

	name := model.Metadata.ServiceAbbreviation
	if len(name) == 0 {
		name = model.Metadata.ServiceFullName
	}
	name = strings.Title(name)

	name = strings.Replace(name, "Amazon", "", -1)
	name = strings.Replace(name, "AWS", "", -1)
	name = serviceIDRegex.ReplaceAllString(name, "")
	name = prefixDigitRegex.ReplaceAllString(name, "")
	name = strings.TrimSpace(name)
	return name
}
