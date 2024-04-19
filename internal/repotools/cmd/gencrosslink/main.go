// package main is the cmd
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type shapeType string

const (
	shapeTypeService shapeType = "service"
	// ...
)

type smithyModel struct {
	Shapes map[string]smithyShape `json:"shapes"`
}

func (m *smithyModel) Service() *serviceShape {
	for _, s := range m.Shapes {
		if s.typ == shapeTypeService {
			return s.asService()
		}
	}
	panic("no service shape found")
}

type smithyShape struct {
	typ shapeType
	raw []byte
}

func (s *smithyShape) UnmarshalJSON(p []byte) error {
	var head = struct {
		Type shapeType `json:"type"`
	}{}
	if err := json.Unmarshal(p, &head); err != nil {
		return nil
	}

	s.typ = head.Type
	s.raw = p
	return nil
}

func (s *smithyShape) asService() *serviceShape {
	var shape serviceShape
	if err := json.Unmarshal(s.raw, &shape); err != nil {
		panic(err)
	}

	return &shape
}

type serviceShape struct {
	Version string       `json:"version"`
	Traits  smithyTraits `json:"traits"`
}

type smithyTraits map[string]json.RawMessage

func (ts smithyTraits) ServiceTrait() (*serviceTrait, bool) {
	const traitID = "aws.api#service"

	raw, ok := ts[traitID]
	if !ok {
		return nil, false
	}

	var v serviceTrait
	if err := json.Unmarshal(raw, &v); err != nil {
		panic(err)
	}

	return &v, true
}

type serviceTrait struct {
	SdkID string `json:"sdkId"`
	DocID string `json:"docId"`
}

func loadModels(dir string) ([]smithyModel, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read models dir: %w", err)
	}

	var models []smithyModel
	for _, file := range files {
		f, err := os.Open(path.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", file.Name(), err)
		}

		p, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file.Name(), err)
		}

		var model smithyModel
		if err := json.Unmarshal(p, &model); err != nil {
			return nil, fmt.Errorf("unmarshal %s: %w", file.Name(), err)
		}

		models = append(models, model)
	}
	return models, nil
}

func mtoc(model *smithyModel) (docID, pkgName string) {
	service := model.Service()
	serviceTrait, ok := service.Traits.ServiceTrait()
	if !ok {
		panic("no service trait")
	}

	pkgName = topkg(serviceTrait.SdkID)
	docID = serviceTrait.DocID
	if docID == "" {
		docID = fmt.Sprintf("%s-%s", normalize2(serviceTrait.SdkID), service.Version)
	}
	return
}

func topkg(sdkID string) (pkg string) {
	pkg = strings.ToLower(sdkID)
	pkg = strings.Replace(pkg, " ", "", -1)
	pkg = strings.Replace(pkg, "-", "", -1)
	return
}

func normalize2(sdkID string) (v string) {
	v = strings.ToLower(sdkID)
	v = strings.Replace(v, " ", "-", -1)
	return
}

func main() {
	var modelsDir string
	flag.StringVar(&modelsDir, "modelsDir", "", "full path to aws-models directory")
	flag.Parse()
	if modelsDir == "" {
		panic("modelsDir not set")
	}

	models, err := loadModels(modelsDir)
	if err != nil {
		panic(err)
	}

	docToPkg := map[string]string{}
	for _, model := range models {
		doc, pkg := mtoc(&model)
		docToPkg[doc] = pkg
	}

	var mapping string
	for k, v := range docToPkg {
		mapping += fmt.Sprintf("            '%s': '%s',\n", k, v)
	}

	page := `<!DOCTYPE html>
<html>
    <script type="text/javascript">
		(function() {
        var docToPkg = {
` + mapping + `        };

        var query = new URL(window.location.href).searchParams;
        var docId = query.get('doc');
        var operation = query.get('operation');
        if (docId === null || operation === null) {
            window.location = 'https://pkg.go.dev/github.com/aws/aws-sdk-go-v2';
            return;
        }

        var service = docToPkg[docId];
        if (service === undefined) {
            window.location = 'https://pkg.go.dev/github.com/aws/aws-sdk-go-v2';
            return;
        }

		window.location = 'https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/' + service + '#Client.' + operation;
		})();
    </script>
</html>`

	fmt.Println(page)
}
