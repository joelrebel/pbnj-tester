package internal

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestName identifies a test
type TestName string

const (
	TestPowerStatus TestName = "power-status"
	TestPowerOn     TestName = "power-on"
	TestPowerOff    TestName = "power-off"
	TestPowerCycle  TestName = "power-cycle"
	TestPxeBoot     TestName = "pxeBoot"
)

// Tester runs tests on a host, this struct holds config attributes for tester
type Tester struct {
	pbnjAddr string
	bmcHost  string
	bmcUser  string
	bmcPass  string
	bmcPort  string
	tests    []test
	results  []Result
	logger   logr.Logger
}

// TestFunc should return output if any and an error to indicate test failure.
type testFunc func(context.Context, *grpc.ClientConn) ([]byte, error)

var (
	// registry of test functions
	registry = map[TestName]testFunc{}

	// SupportedActions should list all actions supported by Tester.
	SupportedActions = []string{
		string(TestPowerStatus),
		string(TestPowerOn),
		string(TestPowerOff),
		string(TestPowerCycle),
		string(TestPxeBoot),
	}
)

func init() {

}

// test holds attributes for a test executed by Tester
type test struct {
	TestName   TestName
	TestMethod testFunc
}

// NewTester returns a Tester instance with the parameters configured.
func NewTester(pbnjAddr, bmcHost, bmcUser, bmcPass, bmcPort string, logLevel string) *Tester {
	return &Tester{
		pbnjAddr: pbnjAddr,
		bmcHost:  bmcHost,
		bmcUser:  bmcUser,
		bmcPass:  bmcPass,
		bmcPort:  bmcPort,
		logger:   NewLogger(logLevel),
	}
}

// Run runs all the given tests.
func (t *Tester) Run(ctx context.Context, testNames []string) {
	if err := t.initTester(ctx, testNames); err != nil {
		t.logger.Error(err, "tester init error")
	}

	conn, err := t.connectGrpc(ctx)
	if err != nil {
		t.logger.Error(err, "connect to PBnJ failed")
		os.Exit(1)
	}

	defer conn.Close()

	t.logger.V(2).Info("PBnJ connection successful.")

	for _, test := range t.tests {
		result := Result{TestName: string(test.TestName)}

		t.logger.V(2).Info("running test", "testName", test.TestName)

		startTime := time.Now()

		output, err := test.TestMethod(ctx, conn)
		result.Output = output
		if err != nil {
			result.Error = err
			result.Runtime = time.Since(startTime)

			t.logger.V(2).Info("Test failed: ", test.TestName)

			t.results = append(t.results, result)
			continue
		}

		result.Succeeded = true
		t.results = append(t.results, result)

		t.logger.V(2).Info("Test successful: ", test.TestName)
	}
}

func (t *Tester) initTester(ctx context.Context, testNames []string) error {
	registry = map[TestName]testFunc{
		TestPowerStatus: t.powerStatus,
		TestPowerOn:     t.powerOn,
		TestPowerOff:    t.powerOff,
		TestPowerCycle:  t.powerCycle,
		TestPxeBoot:     t.pxeBoot,
	}

	// init tests to run
	tests := make([]test, 0, len(testNames))
	for _, testName := range testNames {
		name := TestName(testName)

		f, exists := registry[name]
		if !exists {
			return errors.New("unknown test name: " + testName)
		}

		tests = append(tests, test{
			TestName:   name,
			TestMethod: f,
		})
	}

	t.tests = tests

	return nil
}

func (t *Tester) Results() []Result {
	return t.results
}

// DeviceResult holds the test results for a given device
type DeviceResult struct {
	Vendor  string
	Model   string
	Name    string
	BMCIP   string
	Results []Result
}

// Result is a single test result
type Result struct {
	TestName  string
	Output    []byte
	Error     error
	Succeeded bool
	Runtime   time.Duration
}

// ResultStore stores test results
type ResultStore struct {
	mu      *sync.RWMutex
	results []DeviceResult
}

func NewTestResultStore() *ResultStore {
	return &ResultStore{mu: &sync.RWMutex{}, results: []DeviceResult{}}
}

func (r *ResultStore) Save(result DeviceResult) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.results = append(r.results, result)
}

func (r *ResultStore) Read() []DeviceResult {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.results
}

// NewLogger returns a logr.
func NewLogger(level string) logr.Logger {
	logger := zerolog.New(os.Stdout)

	logger = logger.With().Caller().Timestamp().Logger()

	var l zerolog.Level
	switch level {
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	logger = logger.Level(l)

	return zerologr.New(&logger)
}

func (t *Tester) connectGrpc(ctx context.Context) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(t.pbnjAddr, opts)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
