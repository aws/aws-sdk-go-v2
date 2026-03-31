package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type benchmarkResult struct {
	ID     string  `json:"id"`
	N      int     `json:"n"`
	Mean   float64 `json:"mean"`
	P50    float64 `json:"p50"`
	P90    float64 `json:"p90"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
	StdDev float64 `json:"std_dev"`
}

type benchmarkOutput struct {
	Metadata       benchmarkMetadata `json:"metadata"`
	SerdBenchmarks []benchmarkResult `json:"serde_benchmarks"`
}

type benchmarkMetadata struct {
	Lang      string     `json:"lang"`
	Software  [][]string `json:"software"`
	OS        string     `json:"os"`
	Instance  string     `json:"instance"`
	Precision string     `json:"precision"`
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	var merged benchmarkOutput
	first := true

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || info.Name() != "benchmark.json" {
			return err
		}
		fmt.Println("scanning ", path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var b benchmarkOutput
		if err := json.Unmarshal(data, &b); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}
		if first {
			merged.Metadata = b.Metadata
			first = false
		}
		merged.SerdBenchmarks = append(merged.SerdBenchmarks, b.SerdBenchmarks...)
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	out, _ := json.MarshalIndent(merged, "", "  ")
	fmt.Println(string(out))
}
