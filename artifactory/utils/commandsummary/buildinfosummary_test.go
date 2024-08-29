package commandsummary

import (
	buildinfo "github.com/jfrog/build-info-go/entities"
	"github.com/jfrog/jfrog-cli-core/v2/utils/coreutils"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	buildInfoTable        = "build-info-table.md"
	dockerImageModule     = "docker-image-module.md"
	genericModule         = "generic-module.md"
	mavenModule           = "maven-module.md"
	dockerMultiArchModule = "multiarch-docker-image.md"
)

type MockScanResult struct {
	Violations      string
	Vulnerabilities string
}

// GetViolations returns the mock violations
func (m *MockScanResult) GetViolations() string {
	return m.Violations
}

// GetVulnerabilities returns the mock vulnerabilities
func (m *MockScanResult) GetVulnerabilities() string {
	return m.Vulnerabilities
}

func prepareBuildInfoTest() (*BuildInfoSummary, func()) {
	// Mock the scan results defaults
	StaticMarkdownConfig.scanResultsMapping = make(map[string]ScanResult)
	StaticMarkdownConfig.scanResultsMapping[NonScannedResult] = &MockScanResult{
		Violations:      "Not scanned",
		Vulnerabilities: "Not scanned",
	}
	// Mock config
	StaticMarkdownConfig.setPlatformUrl(testPlatformUrl)
	StaticMarkdownConfig.setPlatformMajorVersion(7)
	StaticMarkdownConfig.setExtendedSummary(false)
	// Cleanup config
	cleanup := func() {
		StaticMarkdownConfig.setExtendedSummary(false)
		StaticMarkdownConfig.setPlatformMajorVersion(0)
		StaticMarkdownConfig.setPlatformUrl("")
	}
	// Create build info instance
	buildInfoSummary := &BuildInfoSummary{}
	return buildInfoSummary, cleanup
}

func TestBuildInfoTable(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:     "buildName",
			Number:   "123",
			Started:  "2024-05-05T12:47:20.803+0300",
			BuildUrl: "http://myJFrogPlatform/builds/buildName/123",
		},
	}
	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoTable(builds)
		testMarkdownOutput(t, getTestDataFile(t, buildInfoTable), res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoTable(builds)
		testMarkdownOutput(t, getTestDataFile(t, buildInfoTable), res)
	})
}

func TestBuildInfoModulesMaven(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:     "buildName",
			Number:   "123",
			Started:  "2024-05-05T12:47:20.803+0300",
			BuildUrl: "http://myJFrogPlatform/builds/buildName/123",
			Modules: []buildinfo.Module{
				{
					Id:   "maven",
					Type: buildinfo.Maven,
					Artifacts: []buildinfo.Artifact{{
						Name:                   "artifact1",
						Path:                   "path/to/artifact1",
						OriginalDeploymentRepo: "libs-release",
					}},
					Dependencies: []buildinfo.Dependency{{
						Id: "dep1",
					}},
				},
			},
		},
	}

	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, mavenModule), res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, mavenModule), res)
	})
}

func TestBuildInfoModulesGradle(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:     "buildName",
			Number:   "123",
			Started:  "2024-05-05T12:47:20.803+0300",
			BuildUrl: "http://myJFrogPlatform/builds/buildName/123",
			Modules: []buildinfo.Module{
				{
					Id:   "gradle",
					Type: buildinfo.Gradle,
					Artifacts: []buildinfo.Artifact{
						{
							Name:                   "gradleArtifact",
							Path:                   "dir/gradleArtifact",
							OriginalDeploymentRepo: "gradle-local",
						},
					},
				},
			},
		},
	}

	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoModules(builds)
		assert.Empty(t, res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoModules(builds)
		assert.Empty(t, res)
	})
}

func TestBuildInfoModulesGeneric(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:     "buildName",
			Number:   "123",
			Started:  "2024-05-05T12:47:20.803+0300",
			BuildUrl: "http://myJFrogPlatform/builds/buildName/123",
			Modules: []buildinfo.Module{
				{
					Id:   "generic",
					Type: buildinfo.Generic,
					Artifacts: []buildinfo.Artifact{{
						Name:                   "artifact2",
						Path:                   "path/to/artifact2",
						OriginalDeploymentRepo: "generic-local",
					}},
				},
			},
		},
	}

	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, genericModule), res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, genericModule), res)
	})
}

func TestDockerModule(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:    "dockerx",
			Number:  "1",
			Started: "2024-08-12T11:11:50.198+0300",
			Modules: []buildinfo.Module{
				{
					Properties: map[string]interface{}{
						"docker.image.tag": "ecosysjfrog.jfrog.io/docker-local/multiarch-image:1",
					},
					Type:   "docker",
					Parent: "image:2",
					Id:     "image:2",
					Checksum: buildinfo.Checksum{
						Sha256: "aae9",
					},
					Artifacts: []buildinfo.Artifact{
						{
							Checksum: buildinfo.Checksum{
								Sha1:   "32c1416f8430fbbabd82cb014c5e09c5fe702404",
								Sha256: "aae9",
								Md5:    "f568bfb1c9576a1f06235ebe0389d2d8",
							},
							Name:                   "sha256__aae9",
							Path:                   "image2/sha256:552c/sha256__aae9",
							OriginalDeploymentRepo: "docker-local",
						},
					},
				},
			},
		},
	}

	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, dockerImageModule), res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, dockerImageModule), res)
	})

}

func TestDockerMultiArchModule(t *testing.T) {
	buildInfoSummary, cleanUp := prepareBuildInfoTest()
	defer func() {
		cleanUp()
	}()
	var builds = []*buildinfo.BuildInfo{
		{
			Name:    "dockerx",
			Number:  "1",
			Started: "2024-08-12T11:11:50.198+0300",
			Modules: []buildinfo.Module{
				{
					Properties: map[string]interface{}{
						"docker.image.tag": "ecosysjfrog.jfrog.io/docker-local/multiarch-image:1",
					},
					Type: "docker",
					Id:   "multiarch-image:1",
					Artifacts: []buildinfo.Artifact{
						{
							Type: "json",
							Checksum: buildinfo.Checksum{
								Sha1:   "fa",
								Sha256: "2217",
								Md5:    "ba0",
							},
							Name:                   "list.manifest.json",
							Path:                   "multiarch-image/1/list.manifest.json",
							OriginalDeploymentRepo: "docker-local",
						},
					},
				},
				{
					Type:   "docker",
					Parent: "multiarch-image:1",
					Id:     "linux/amd64/multiarch-image:1",
					Artifacts: []buildinfo.Artifact{
						{
							Checksum: buildinfo.Checksum{
								Sha1:   "32",
								Sha256: "sha256:552c",
								Md5:    "f56",
							},
							Name:                   "manifest.json",
							Path:                   "multiarch-image/sha256",
							OriginalDeploymentRepo: "docker-local",
						},
						{
							Checksum: buildinfo.Checksum{
								Sha1:   "32c",
								Sha256: "aae9",
								Md5:    "f56",
							},
							Name:                   "sha256__aae9",
							Path:                   "multiarch-image/sha256:552c/sha256",
							OriginalDeploymentRepo: "docker-local",
						},
					},
				},
			},
		},
	}

	t.Run("Extended Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(true)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, dockerMultiArchModule), res)
	})
	t.Run("Basic Summary", func(t *testing.T) {
		StaticMarkdownConfig.setExtendedSummary(false)
		res := buildInfoSummary.buildInfoModules(builds)
		testMarkdownOutput(t, getTestDataFile(t, dockerMultiArchModule), res)
	})

}

// Tests data files are location artifactory/commands/testdata/command_summary
func getTestDataFile(t *testing.T, fileName string) string {
	var modulesPath string
	if StaticMarkdownConfig.IsExtendedSummary() {
		modulesPath = filepath.Join("../", "testdata", "command_summaries", "extended", fileName)
	} else {
		modulesPath = filepath.Join("../", "testdata", "command_summaries", "basic", fileName)
	}

	content, err := os.ReadFile(modulesPath)
	assert.NoError(t, err)
	contentStr := string(content)
	if coreutils.IsWindows() {
		contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	}
	return contentStr
}

// Sometimes there are inconsistencies in the Markdown output, this function normalizes the output for comparison
// This allows easy debugging when tests fails
func normalizeMarkdown(md string) string {
	// Remove the extra spaces added for equal table length
	md = strings.ReplaceAll(md, markdownSpaceFiller, "")
	md = strings.ReplaceAll(md, "\r\n", "\n")
	md = strings.ReplaceAll(md, "\r", "\n")
	md = strings.ReplaceAll(md, `\n`, "\n")
	// Regular expression to match the table rows and header separators
	re := regexp.MustCompile(`\s*\|\s*`)
	// Normalize spaces around the pipes and colons in the Markdown
	lines := strings.Split(md, "\n")
	for i, line := range lines {
		if strings.Contains(line, "|") {
			// Remove extra spaces around pipes and colons
			line = re.ReplaceAllString(line, " | ")
			lines[i] = strings.TrimSpace(line)
		}
	}
	return strings.Join(lines, "\n")
}

func testMarkdownOutput(t *testing.T, expected, actual string) {
	expected = normalizeMarkdown(expected)
	actual = normalizeMarkdown(actual)

	// If the compared string length exceeds the maximum length,
	// the string is not formatted, leading to an unequal comparison.
	// Ensure to test small units of Markdown for better unit testing
	// and to facilitate testing.
	maxCompareLength := 950
	if len(expected) > maxCompareLength || len(actual) > maxCompareLength {
		t.Fatalf("Markdown output is too long to compare, limit the length to %d chars", maxCompareLength)
	}
	assert.Equal(t, expected, actual)
}
