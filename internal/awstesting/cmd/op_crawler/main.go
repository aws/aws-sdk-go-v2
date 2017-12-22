package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func main() {
	var filterService, filterOp string
	flag.StringVar(&filterService, "s", "", "The `service` to execute. If unset will run all services.")
	flag.StringVar(&filterOp, "o", "", "The `operation` to execute. If unset will run all operations. Requires service set.")
	flag.Parse()

	if len(filterOp) > 0 && len(filterService) == 0 {
		flag.PrintDefaults()
		panic("operation filter requires service set also")
	}

	w := writer{
		buf:       bytes.NewBuffer(nil),
		indentStr: "\t",
	}
	server := setupServer(w.Indent())
	defer server.Close()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("failed to load config " + err.Error())
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL)
	cfg.Region = endpoints.UsWest2RegionID
	cfg.Credentials = aws.AnonymousCredentials
	cfg.Handlers.Validate.Remove(defaults.ValidateParametersHandler)

	cfg.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	cfg.Handlers.Send.PushFront(func(r *aws.Request) {
		w.Writef("%s.%s:", r.Metadata.ServiceName, r.Operation.Name)
	})

	for _, service := range createServices(cfg) {
		if len(filterService) > 0 && service.name != filterService {
			continue
		}

		fmt.Println("Processing:", service.name)
		if err := callService(service.value, filterOp); err != nil {
			panic(fmt.Sprintf("Service, %s failed, %v", service.name, err))
		}
	}

	io.Copy(os.Stdout, w.buf)
}

const (
	allowedRecursion = 2
	sliceSize        = 3
)

func callService(svcV reflect.Value, filterOp string) error {
	reqOption := aws.Option(
		func(r *aws.Request) {
			r.Handlers.Build.PushBack(func(req *aws.Request) {
				endpoint, _ := r.Config.EndpointResolver.ResolveEndpoint(r.Metadata.ServiceName, r.Config.Region)
				origURL, _ := url.Parse(endpoint.URL)
				// Correct ContentLength fields for S3 operations
				switch r.Metadata.ServiceName {
				case "s3":
					switch r.Operation.Name {
					case "PutObject", "UploadPart":
						n, _ := computeBodyLength(r.GetBody())
						r.HTTPRequest.Header.Set("Content-Length", strconv.FormatInt(n, 10))
					}
				case "machinelearning":
					newURL := r.HTTPRequest.URL
					r.HTTPRequest.URL = origURL
					r.HTTPRequest.URL.Path = newURL.Path
					if !strings.HasPrefix(r.HTTPRequest.URL.Path, "/") {
						r.HTTPRequest.URL.Path = "/" + r.HTTPRequest.URL.Path
					}
					r.HTTPRequest.URL.RawPath = newURL.RawPath
					r.HTTPRequest.URL.RawQuery = newURL.RawQuery
				}
			})

			r.Handlers.Complete.PushBack(func(req *aws.Request) {
				if r.Error != nil {
					fmt.Println(r.Params)
				}
			})
		},
	)

	svcT := svcV.Type()
	n := svcT.NumMethod()

	ops := []string{}
	for i := 0; i < n; i++ {
		fm := svcT.Method(i)
		fName := fm.Name

		if fName == "NewRequest" || !strings.HasSuffix(fName, "Request") || (len(filterOp) > 0 && fName != filterOp) {
			continue
		}

		ops = append(ops, strings.TrimSuffix(fName, "Request"))
	}

	sort.Strings(ops)
	for _, op := range ops {
		fName := op + "Request"

		fv := svcV.MethodByName(fName)
		fm, _ := svcT.MethodByName(fName)
		ft := fm.Type

		fmt.Println("-", op)

		it := ft.In(1)
		iv := valueForType(it, visitType(it))

		ovs := fv.Call([]reflect.Value{iv})

		m := findMethod(ovs[0], "ApplyOptions")
		m.Call([]reflect.Value{reflect.ValueOf(reqOption)})

		m = findMethod(ovs[0], "Send")
		ovs = m.Call([]reflect.Value{})

		if v := ovs[1]; !v.IsNil() {
			return v.Interface().(error)
		}
	}

	return nil
}

// findMethod finds the method name on the v type using a case-insensitive
// lookup. Returns nil if no method is found.
func findMethod(v reflect.Value, methodName string) *reflect.Value {
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		if name := t.Method(i).Name; name == methodName {
			m := v.MethodByName(name)
			return &m
		}
	}

	// Check if the method exists on a field
	for i := 0; i < t.NumField(); i++ {
		if m := findMethod(reflect.Indirect(v).Field(i), methodName); m != nil {
			return m
		}
	}

	return nil
}

func asVisited(v map[reflect.Type]int, t reflect.Type) (map[reflect.Type]int, bool) {
	if c, ok := v[t]; ok {
		if c == 0 {
			return v, false
		}
		c--
		v[t] = c
	} else {
		v[t] = allowedRecursion
	}

	return v, true
}

var ioReadSeekerType = reflect.TypeOf((*io.ReadSeeker)(nil)).Elem()
var emptyInterfaceType = reflect.TypeOf((*interface{})(nil)).Elem()

func valueForType(vt reflect.Type, visited *visitedType) reflect.Value {
	var v reflect.Value

	vtt := vt
	if vtt.Kind() == reflect.Ptr {
		vtt = vtt.Elem()
	}

	switch vtt.Kind() {
	case reflect.Map:
		v = reflect.MakeMap(vtt)
		kt := vtt.Key()
		kv := valueForType(kt, visited)
		et := vtt.Elem()
		ev := valueForType(et, visited)
		v.SetMapIndex(kv, ev)

	case reflect.Slice:
		v = reflect.MakeSlice(vtt, sliceSize, sliceSize)
		vet := vtt.Elem()
		for i := 0; i < sliceSize; i++ {
			sv := v.Index(i)
			nsv := valueForType(vet, visited)
			sv.Set(nsv)
		}

	case reflect.Interface:
		switch vtt {
		case ioReadSeekerType:
			v = reflect.New(vtt)
			v.Elem().Set(reflect.ValueOf(bytes.NewReader([]byte("byte value"))))
		case emptyInterfaceType:
			v = reflect.ValueOf("empty interface value")
		default:
			panic("value for interface, unknown type" + vtt.String())
		}

	case reflect.String:
		v = reflect.New(vtt)
		v.Elem().SetString("stringValue")

	case reflect.Bool:
		v = reflect.New(vtt)
		v.Elem().SetBool(true)

	case reflect.Uint8: // byte
		v = reflect.New(vtt)
		v.Elem().Set(reflect.ValueOf(uint8('b')))

	case reflect.Int64:
		v = reflect.New(vtt)
		v.Elem().SetInt(987654321)

	case reflect.Float64:
		v = reflect.New(vtt)
		v.Elem().SetFloat(123456789.321)

	case reflect.Struct:
		v = reflect.New(vtt)
		ve := v.Elem()
		n := ve.NumField()
		for i := 0; i < n; i++ {
			fv := ve.Field(i)
			fs := vtt.Field(i)
			ft := fv.Type()
			if len(fs.PkgPath) != 0 {
				continue
			}
			nested, keep := visited.Visit(ft)
			if !keep {
				continue
			}
			nfv := valueForType(ft, nested)
			fv.Set(nfv)
		}
	default:
		panic("unknown type, " + vtt.String())
	}

	if vt.Kind() != reflect.Ptr && v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}

func setupServer(out writer) *httptest.Server {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		out.Writef("Path:")
		out.Indent().Writef(r.URL.Path)
		out.Writef("Query:")
		var keys []string
		query := r.URL.Query()
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			for _, v := range query[k] {
				out.Indent().Writef("%s: %s", k, v)
			}
		}
		out.Writef("Headers:")
		keys = keys[0:0]
		for k := range r.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			for _, v := range r.Header[k] {
				out.Indent().Writef("%s: %s", k, v)
			}
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		r.Body.Close()
		out.Writef("Body:")
		out.Indent().Writef(string(body))
	}))

	return server
}

type writer struct {
	buf       *bytes.Buffer
	indent    int
	indentStr string
}

func (w writer) Writef(format string, args ...interface{}) error {
	indent := strings.Repeat(w.indentStr, w.indent)
	w.buf.WriteString(indent)
	w.buf.WriteString(fmt.Sprintf(format, args...))
	w.buf.WriteRune('\n')
	return nil
}
func (w writer) Indent() writer {
	newW := w
	newW.indent++
	return newW
}

type visitedType struct {
	typ  reflect.Type
	left int
	next *visitedType
}

func visitType(t reflect.Type) *visitedType {
	return &visitedType{
		typ:  t,
		left: allowedRecursion,
	}
}

func (v *visitedType) String() string {
	if v == nil {
		return "END"
	}
	return fmt.Sprintf("Type:%v,Kind:%v,Left:%v->%v", v.typ.Name(), v.typ.Kind(), v.left, v.next.String())
}

func (v *visitedType) copy() *visitedType {
	nv := &visitedType{}

	oldNext := v
	newNext := nv

	for {
		*newNext = *oldNext
		oldNext = oldNext.next

		if oldNext == nil {
			break
		}

		newNext.next = &visitedType{}
		newNext = newNext.next
	}

	return nv
}

func (v *visitedType) Visit(t reflect.Type) (*visitedType, bool) {
	if v == nil {
		return visitType(t), true
	}

	nv := v.copy()

	last := nv
	next := nv
	for next != nil {
		last = next
		if next.typ != t {
			next = next.next
			continue
		}

		next.left--
		return nv, next.left >= 0
	}

	last.next = visitType(t)

	return nv, true
}

func computeBodyLength(r io.ReadSeeker) (int64, error) {
	seekable := true
	// Determine if the seeker is actually seekable. ReaderSeekerCloser
	// hides the fact that a io.Readers might not actually be seekable.
	switch v := r.(type) {
	case aws.ReaderSeekerCloser:
		seekable = v.IsSeeker()
	case *aws.ReaderSeekerCloser:
		seekable = v.IsSeeker()
	}
	if !seekable {
		return -1, nil
	}

	curOffset, err := r.Seek(0, 1)
	if err != nil {
		return 0, err
	}

	endOffset, err := r.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	_, err = r.Seek(curOffset, 0)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}
