package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	modelPath  string
	outputPath string
)

func init() {
	flag.StringVar(&modelPath, "m", "", "The `path` of the source API models to search for smithy files.")
	flag.StringVar(&outputPath, "o", "", "The `path` the API models are copied to")
}

func main() {
	flag.Parse()
	if len(modelPath) == 0 || len(outputPath) == 0 {
		log.Fatalf("model path and output path required")
		flag.PrintDefaults()
	}

	srcModels, err := findSmithyModels(modelPath)
	if err != nil {
		log.Fatalf("failed to get models, %v", err)
	}

	for _, model := range srcModels {
		if err = copyModelFile(model); err != nil {
			log.Fatalf("copy failed, %v", err)
		}
	}
}

func copyModelFile(model SourceModel) error {
	srcFile, err := os.Open(model.SrcFilepath)
	if err != nil {
		return fmt.Errorf("failed to open source file %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(filepath.Join(outputPath, model.DstFilename),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create destination file %w", err)
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy model, %w", err)
	}

	if err = dstFile.Close(); err != nil {
		return fmt.Errorf("failed to close destination file, %w", err)
	}

	return nil
}

// SourceModel provides the type for a model that should be copied.
type SourceModel struct {
	SDKID       string
	Version     string
	SrcFilepath string
	DstFilename string
}

func findSmithyModels(modelPath string) (map[string]SourceModel, error) {
	models := map[string]SourceModel{}

	err := filepath.Walk(modelPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			if filepath.Ext(info.Name()) != ".json" {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			var check SmithyCheck
			if err = json.NewDecoder(f).Decode(&check); err != nil {
				log.Printf("skipping %s file, json but unable to decode, %v",
					path, err)
				return err
			}
			if !strings.HasPrefix(check.Smithy, "1.") && !strings.HasPrefix(check.Smithy, "2.") {
				return nil
			}
			var shapes map[string]SmithyShape
			if err = json.Unmarshal(check.Shapes, &shapes); err != nil {
				return fmt.Errorf("failed to unmarshal smithy model %v, %w",
					path, err)
			}

			for name, shape := range shapes {
				if shape.Type != "service" {
					continue
				}
				if len(shape.Version) == 0 {
					return fmt.Errorf("smithy service doesn't have version %s %s",
						name, path)
				}
				if shape.Traits.Service == nil {
					// Ignore services that don't have an SDK id.
					continue
				}
				if len(shape.Traits.Service.SDKID) == 0 {
					return fmt.Errorf("smithy service doesn't have sdkId value %s, %s",
						name, path)
				}

				sdkID := shape.Traits.Service.SDKID
				sdkID = strings.ReplaceAll(sdkID, " ", "-")
				sdkID = strings.ToLower(sdkID)
				if o, ok := models[sdkID]; ok {
					return fmt.Errorf("two smithy models have same sdkId %s, 1:%s 2:%s",
						sdkID, o.SrcFilepath, path)
				}
				// TODO what about two services in same model file?
				models[sdkID] = SourceModel{
					SDKID:       sdkID,
					Version:     shape.Version,
					SrcFilepath: path,
					DstFilename: sdkID + ".json",
				}
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return models, nil
}

// SmithyCheck provides initial JSON object deserialization checking for a
// Smithy json model.
type SmithyCheck struct {
	Smithy string          `json:"smithy"`
	Shapes json.RawMessage `json:"shapes"`
}

// SmithyShape provides JSON object deserialization type for a Smithy Shape
type SmithyShape struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Traits  struct {
		Service *struct {
			SDKID string `json:"sdkId"`
		} `json:"aws.api#service"`
	} `json:"traits"`
}
