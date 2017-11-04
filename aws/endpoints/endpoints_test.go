package endpoints

import "testing"

func TestDefaultResolver_Partitions(t *testing.T) {
	resolver := NewDefaultResolver()
	ps := resolver.Partitions()

	if a, e := len(ps), len(defaultPartitions); a != e {
		t.Errorf("expected %d partitions, got %d", e, a)
	}
}

func TestEnumDefaultRegions(t *testing.T) {
	expectPart := defaultPartitions[0]
	partEnum := defaultPartitions[0].Partition()

	regEnum := partEnum.Regions()

	if a, e := len(regEnum), len(expectPart.Regions); a != e {
		t.Errorf("expected %d regions, got %d", e, a)
	}
}

func TestEnumPartitionServices(t *testing.T) {
	expectPart := testPartitions[0]
	partEnum := testPartitions[0].Partition()

	if a, e := partEnum.ID(), "part-id"; a != e {
		t.Errorf("expect %q partition ID, got %q", e, a)
	}

	svcEnum := partEnum.Services()

	if a, e := len(svcEnum), len(expectPart.Services); a != e {
		t.Errorf("expected %d regions, got %d", e, a)
	}
}

func TestEnumRegionServices(t *testing.T) {
	p := testPartitions[0].Partition()

	rs := p.Regions()

	if a, e := len(rs), 2; a != e {
		t.Errorf("expect %d regions got %d", e, a)
	}

	if _, ok := rs["us-east-1"]; !ok {
		t.Errorf("expect us-east-1 region to be found, was not")
	}
	if _, ok := rs["us-west-2"]; !ok {
		t.Errorf("expect us-west-2 region to be found, was not")
	}

	r := rs["us-east-1"]

	if a, e := r.ID(), "us-east-1"; a != e {
		t.Errorf("expect %q region ID, got %q", e, a)
	}

	ss := r.Services()
	if a, e := len(ss), 1; a != e {
		t.Errorf("expect %d services for us-east-1, got %d", e, a)
	}

	if _, ok := ss["service1"]; !ok {
		t.Errorf("expect service1 service to be found, was not")
	}

	resolved, err := r.Endpoint("service1", ResolveOptions{})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if a, e := resolved.URL, "https://service1.us-east-1.amazonaws.com"; a != e {
		t.Errorf("expect %q resolved URL, got %q", e, a)
	}
}

func TestEnumServiceRegions(t *testing.T) {
	p := testPartitions[0].Partition()

	rs := p.Services()["service1"].Regions()
	if e, a := 2, len(rs); e != a {
		t.Errorf("expect %d regions, got %d", e, a)
	}

	if _, ok := rs["us-east-1"]; !ok {
		t.Errorf("expect region to be found")
	}
	if _, ok := rs["us-west-2"]; !ok {
		t.Errorf("expect region to be found")
	}
}

func TestEnumServicesEndpoints(t *testing.T) {
	p := testPartitions[0].Partition()

	ss := p.Services()

	if a, e := len(ss), 5; a != e {
		t.Errorf("expect %d regions got %d", e, a)
	}

	if _, ok := ss["service1"]; !ok {
		t.Errorf("expect service1 region to be found, was not")
	}
	if _, ok := ss["service2"]; !ok {
		t.Errorf("expect service2 region to be found, was not")
	}

	s := ss["service1"]
	if a, e := s.ID(), "service1"; a != e {
		t.Errorf("expect %q service ID, got %q", e, a)
	}

	resolved, err := s.Endpoint("us-west-2", ResolveOptions{})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if a, e := resolved.URL, "https://service1.us-west-2.amazonaws.com"; a != e {
		t.Errorf("expect %q resolved URL, got %q", e, a)
	}
}

func TestEnumEndpoints(t *testing.T) {
	p := testPartitions[0].Partition()
	s := p.Services()["service1"]

	es := s.Endpoints()
	if a, e := len(es), 2; a != e {
		t.Errorf("expect %d endpoints for service2, got %d", e, a)
	}
	if _, ok := es["us-east-1"]; !ok {
		t.Errorf("expect us-east-1 to be found, was not")
	}

	e := es["us-east-1"]
	if a, e := e.ID(), "us-east-1"; a != e {
		t.Errorf("expect %q endpoint ID, got %q", e, a)
	}
	if a, e := e.ServiceID(), "service1"; a != e {
		t.Errorf("expect %q service ID, got %q", e, a)
	}

	resolved, err := e.Resolve(ResolveOptions{})
	if err != nil {
		t.Fatalf("expect no error, got %v", err)
	}

	if a, e := resolved.URL, "https://service1.us-east-1.amazonaws.com"; a != e {
		t.Errorf("expect %q resolved URL, got %q", e, a)
	}
}

func TestResolveEndpointForPartition(t *testing.T) {
	p := testPartitions.Partitions()[0]

	expected, err := testPartitions.EndpointFor("service1", "us-east-1", ResolveOptions{})

	actual, err := p.Endpoint("service1", "us-east-1", ResolveOptions{})
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if expected != actual {
		t.Errorf("expect resolved endpoint to be %v, but got %v", expected, actual)
	}
}

func TestPartition_RegionsForService(t *testing.T) {
	ps := DefaultPartitions()

	var expect map[string]Region
	var serviceID string
	for _, s := range ps[0].Services() {
		expect = s.Regions()
		serviceID = s.ID()
		if len(expect) > 0 {
			break
		}
	}

	p, ok := ps.ForPartition(ps[0].ID())
	if !ok {
		t.Fatalf("expect partition to exist")
	}

	actual, ok := p.RegionsForService(serviceID)
	if !ok {
		t.Fatalf("expect service to exist")
	}
	if len(actual) == 0 {
		t.Fatalf("expect service %s to have regions", serviceID)
	}
	if e, a := len(expect), len(actual); e != a {
		t.Fatalf("expect %d regions, got %d", e, a)
	}

	for id, r := range actual {
		if e, a := id, r.ID(); e != a {
			t.Errorf("expect %s region id, got %s", e, a)
		}
		if _, ok := expect[id]; !ok {
			t.Errorf("expect %s region to be found", id)
		}
	}
}

func TestRegionsForService_NotFound(t *testing.T) {
	ps := testPartitions.Partitions()

	p, ok := ps.ForPartition(ps[0].ID())
	if !ok {
		t.Fatalf("expect partition to exist")
	}

	actual, ok := p.RegionsForService("service-not-exists")
	if ok {
		t.Fatalf("expect service to not exist")
	}
	if len(actual) != 0 {
		t.Errorf("expect no regions, got %v", actual)
	}
}

func TestPartitionForRegion(t *testing.T) {
	ps := DefaultPartitions()
	expect := ps[len(ps)%2]

	var regionID string
	for id := range expect.Regions() {
		regionID = id
		break
	}

	actual, ok := ps.ForRegion(regionID)
	if !ok {
		t.Fatalf("expect partition to be found")
	}
	if e, a := expect.ID(), actual.ID(); e != a {
		t.Errorf("expect %s partition, got %s", e, a)
	}
}

func TestPartitionForRegion_NotFound(t *testing.T) {
	ps := DefaultPartitions()

	actual, ok := ps.ForRegion("regionNotExists")
	if ok {
		t.Errorf("expect no partition to be found, got %v", actual)
	}
}
