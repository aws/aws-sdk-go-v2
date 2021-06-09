package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/internal/repotools"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/gomod"
	"github.com/aws/aws-sdk-go-v2/internal/repotools/manifest"
	"golang.org/x/mod/modfile"
)

const manifestFileName = "generated.json"

var config = struct {
	BuildArtifactPath string
}{}

func init() {
	flag.StringVar(&config.BuildArtifactPath, "build", "", "build artifact path")
}

func main() {
	flag.Parse()

	if len(config.BuildArtifactPath) == 0 {
		log.Fatalf("expect build artifact path to be provided")
	}

	if stat, err := os.Stat(filepath.Join(config.BuildArtifactPath)); err != nil {
		log.Fatalf("failed to stat build artifact path: %v", err)
	} else if !stat.IsDir() {
		log.Fatalf("build artifact path must be a directory")
	}

	repoRoot, err := repotools.GetRepoRoot()
	if err != nil {
		log.Fatalf("failed to get repository root: %v", err)
	}

	rootMod, err := gomod.LoadModuleFile(repoRoot, nil, true)
	if err != nil {
		log.Fatalf("failed to read repo root go module, %v", err)
	}

	rootModulePath, err := gomod.GetModulePath(rootMod)
	if err != nil {
		log.Fatalf("unable to determine repo root module path, %v", err)
	}

	av := manifest.SmithyArtifactPaths{}
	if err = filepath.Walk(config.BuildArtifactPath, av.Walk); err != nil {
		log.Fatalf("failed to walk build artifacts: %v", err)
	}

	if len(av) == 0 {
		log.Printf("no build artifacts found: %v", err)
		return
	}

	if err := copyBuildArtifacts(av, rootModulePath, repoRoot); err != nil {
		log.Fatalf("failed to copy build artifacts: %v", err)
	}
}

func copyBuildArtifacts(paths []string, rootModulePath string, repoRoot string) error {
	for _, artifactPath := range paths {
		buildManifest, err := manifest.LoadManifest(filepath.Join(artifactPath, manifestFileName))
		if err != nil {
			return fmt.Errorf("failed to load manifest: %w", err)
		}
		if !strings.HasPrefix(buildManifest.Module, rootModulePath) {
			return fmt.Errorf("%v is not a sub-module of %v", buildManifest.Module, rootModulePath)
		}

		moduleRelativePath := strings.TrimPrefix(strings.TrimPrefix(buildManifest.Module, rootModulePath), "/")
		if moduleRelativePath == "" {
			moduleRelativePath = "."
		}

		targetPath := filepath.Join(repoRoot, moduleRelativePath)
		if err := prepareTargetDirectory(targetPath, buildManifest); err != nil {
			return fmt.Errorf("failed to prepare target directory: %w", err)
		}

		if err := copyBuildArtifactToPath(artifactPath, targetPath, buildManifest); err != nil {
			return fmt.Errorf("failed to copy build artifact to target: %w", err)
		}

		generated, err := generateModuleDefinition(buildManifest)
		if err != nil {
			return fmt.Errorf("failed to generate go module file: %w", err)
		}

		err = gomod.WriteModuleFile(targetPath, generated)
		if err != nil {
			return fmt.Errorf("failed to write go module file: %w", err)
		}
	}
	return nil
}

func generateModuleDefinition(m manifest.Manifest) (*modfile.File, error) {
	mod := modfile.File{
		Syntax: &modfile.FileSyntax{},
	}

	if err := mod.AddModuleStmt(m.Module); err != nil {
		return nil, fmt.Errorf("failed to set module path: %v", err)
	}

	if err := mod.AddGoStmt(m.Go); err != nil {
		return nil, fmt.Errorf("failed to set Go version: %v", err)
	}

	for depPath, depVersion := range m.Dependencies {
		depPath := path.Clean(depPath)

		if err := mod.AddRequire(depPath, depVersion); err != nil {
			return nil, fmt.Errorf("failed to add dependency %v@%v", depPath, depVersion)
		}
	}

	return &mod, nil
}

func prepareTargetDirectory(path string, buildManifest manifest.Manifest) error {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	} else if err != nil {
		return err
	}

	targetManifest, err := manifest.LoadManifest(filepath.Join(path, manifestFileName))
	var notFoundErr *manifest.NoManifestFound
	if err != nil && !errors.As(err, &notFoundErr) {
		return err
	}

	var (
		targetModule string
		cleanupList  []string
	)
	if err == nil {
		targetModule = targetManifest.Module
		cleanupList = targetManifest.Files
	} else {
		log.Printf("[WARN] target directory %v is missing generated.json, will only remove files present in build artifact", path)
		if ok, err := gomod.IsGoModPresent(path); err != nil {
			return err
		} else if ok {
			moduleFile, err := gomod.LoadModuleFile(path, nil, true)
			if err != nil {
				return err
			}
			targetModule, err = gomod.GetModulePath(moduleFile)
			if err != nil {
				return err
			}

			if targetModule != buildManifest.Module {
				return fmt.Errorf("target module %v does not match build artifact %v", targetModule, buildManifest.Module)
			}
		}
		cleanupList = buildManifest.Files
	}

	for _, fileName := range cleanupList {
		filePath := filepath.Join(path, fileName)
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %v: %w", filePath, err)
		}
	}

	return nil
}

func copyBuildArtifactToPath(source, target string, m manifest.Manifest) error {
	for _, fp := range m.Files {
		sfp := filepath.Join(source, fp)
		tfp := filepath.Join(target, fp)

		if err := copyArtifact(sfp, tfp); err != nil {
			return err
		}
	}
	return nil
}

func copyArtifact(sourcePath, targetPath string) (err error) {
	dirs, _ := filepath.Split(targetPath)
	if len(dirs) != 0 {
		err = os.MkdirAll(dirs, 0755)
		if err != nil {
			return err
		}
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		fErr := targetFile.Close()
		if fErr != nil && err == nil {
			err = fErr
		}
	}()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}
