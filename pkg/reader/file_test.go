package reader

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/Juli3nnicolas/http_log_monitor/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestOpenDoesntReturnErrorOnExistingFile(t *testing.T) {
	r := File{}
	err := r.Open("./data.test")
	defer r.Close()
	assert.Nil(t, err)
}

func TestOpenReturnsAnErrorIfFileDoesntExist(t *testing.T) {
	r := File{}
	err := r.Open("./foo.bar")
	defer r.Close()
	assert.NotNil(t, err)
}

func TestOpenReturnsErrorIfNotExactlyOneStringParameterIsPassed(t *testing.T) {
	r := File{}

	// Check first param must be a string
	err := r.Open(1)
	defer r.Close()
	assert.NotNil(t, err)

	// Check 2 string params triggers an error
	r = File{}
	err = r.Open("hello", "world")
	defer r.Close()
	assert.NotNil(t, err)
}

func TestCloseProperlyResetPointers(t *testing.T) {
	// Setup
	r := File{}
	err := r.Open("./data.test")
	defer r.Close()
	assert.Nil(t, err)
	assert.NotNil(t, r.file)
	assert.NotNil(t, r.scanner)

	// Exercise
	r.Close()

	// Validation
	assert.Nil(t, r.file)
	assert.Nil(t, r.scanner)
}

func TestCloseCanBecalledSeveralTimesInARowWithoutPanicking(t *testing.T) {
	r := File{}
	err := r.Open("./data.test")
	defer r.Close()
	assert.Nil(t, err)

	r.Close()
	r.Close()
	r.Close()
	r.Close()
}

func TestReadOutputAllLines(t *testing.T) {
	// Setup
	r := File{}
	r.Parse = func(data []byte) (log.Info, error) {
		return log.Info{Host: string(data)}, nil
	}

	// Read the file using the std library to compare
	f, err := os.Open("./data.test")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	expectedResp := make([]log.Info, 0, len(lines))
	for _, l := range lines {
		expectedResp = append(expectedResp, log.Info{Host: l})
	}

	// Exercise
	err = r.Open("./data.test")
	defer r.Close()
	assert.Nil(t, err)

	responses := make([]log.Info, 0, len(expectedResp))
	resp, err := r.Read()

	for resp != nil && err == nil {
		assert.Len(t, resp, 1)
		responses = append(responses, resp[0])
		resp, err = r.Read()
	}
	assert.Nil(t, err)

	// Validation
	assert.Equal(t, expectedResp, responses)
}
