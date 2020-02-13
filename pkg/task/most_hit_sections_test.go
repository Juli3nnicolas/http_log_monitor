package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractSectionReturnsSlashWhenGivenRootPath(t *testing.T) {
	section := extractSection("/")
	assert.Equal(t, "/", section)
}

func TestExtractSectionReturnsRightValueWhenGivenLongPath(t *testing.T) {
	section := extractSection("/mywebsite/js/pkg/foo/script.js")
	assert.Equal(t, "/mywebsite", section)
}

func TestExtractSectionReturnsRightValueWhenGivenShortestAllowedPath(t *testing.T) {
	section := extractSection("/instance/create")
	assert.Equal(t, "/instance", section)
}

func TestExtractSectionReturnsSlashWhenGivenAPathWithASingleSlash(t *testing.T) {
	section := extractSection("/index.html")
	assert.Equal(t, "/", section)
}

func TestExtractSectionReturnsAnEmptyStringWhenGivenAStringWithoutSlashes(t *testing.T) {
	section := extractSection("index.html")
	assert.Equal(t, "", section)
}

func TestExtractSectionReturnsAnEmptyStringWhenGivenAnEmptyString(t *testing.T) {
	section := extractSection("index.html")
	assert.Equal(t, "", section)
}
