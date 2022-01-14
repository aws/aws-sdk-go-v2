package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// Import is a Go package import and associated alias
type Import struct {
	Package string
	Alias   string
}

type generationContext struct {
	PackageName  string
	ResolverName string
	Imports      map[Import]struct{}
	Config       SDKDefaultConfig
}

func (g *generationContext) AddImport(pkg, alias string) (s string) {
	if g.Imports == nil {
		g.Imports = make(map[Import]struct{})
	}

	g.Imports[Import{
		Package: pkg,
		Alias:   alias,
	}] = struct{}{}

	return s
}

func (g *generationContext) AddSDKImport(pkg, alias string) (s string) {
	if g.Imports == nil {
		g.Imports = make(map[Import]struct{})
	}

	g.Imports[Import{
		Package: path.Join("github.com/aws/aws-sdk-go-v2", pkg),
		Alias:   alias,
	}] = struct{}{}

	return s
}

func generateConfigPackage(jsonFile, outputFile, packageName, resolverName string) (err error) {
	config, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	decoder := newDecoder(config)

	var schemaVersion SchemaVersion
	if err = decoder.Decode(&schemaVersion); err != nil {
		return fmt.Errorf("failed to get schema version: %w", err)
	}

	decoder = newDecoder(config)

	if schemaVersion.Version != 1 {
		return fmt.Errorf("generator only supports version 1 schema, got %d", schemaVersion.Version)
	}

	var defaultConfig SDKDefaultConfig
	if err := decoder.Decode(&defaultConfig); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	oFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cErr := oFile.Close(); cErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cErr)
		}
	}()

	configContent := bytes.NewBuffer(nil)

	data := &generationContext{
		PackageName:  packageName,
		ResolverName: resolverName,
		Config:       defaultConfig,
	}

	if err := tmpl.ExecuteTemplate(configContent, "config", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := tmpl.ExecuteTemplate(oFile, "header", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if _, err := io.Copy(oFile, configContent); err != nil {
		return fmt.Errorf("failed to copy to output file: %w", err)
	}

	return nil
}

func newDecoder(data []byte) *json.Decoder {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder
}
