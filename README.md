# feat-serde2-benchmark

The benchmark models have been dropped into this repo and clients have been generated for them. They live in `service/serdbenchmark`.

## Running a test

```
# check out the SDK and the repo to the same parent directory
git clone -b feat-serde2-benchmark git@github.com:aws/aws-sdk-go-v2.git
git clone -b feat-serde2-benchmark git@github.com:aws/smithy-go.get

# point the benchmark services at the local copy of smithy-go
cd aws-sdk-go-v2
make gen-mod-replace-smithy-service_serdbenchmark

# cd into one of the service dirs and test
# results are written to benchmark.json
cd service/serdbenchmark/jsonrpc10dataplane
go test ./...
```
