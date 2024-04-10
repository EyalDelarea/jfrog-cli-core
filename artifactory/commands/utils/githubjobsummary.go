package utils

import (
	"fmt"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"os"
	"path"
)

type Operation string

const (
	Upload               Operation = "rt upload"
	Publish              Operation = "publish"
	GithubEnvStepSummary           = "GITHUB_STEP_SUMMARY"
)

func (o Operation) String() string {
	return "Command: " + string(o)
}

type MarkdownGenerator struct {
	file           *os.File
	result         *Result
	operationTitle Operation
}

func GenerateSummaryMarkdown(result *Result, operationTitle Operation) error {
	githubMarkdownGenerator, cleanUp, err := NewGithubMarkdownGenerator(result, operationTitle)
	defer func() {
		err = cleanUp()
	}()
	if err != nil {
		return err
	}
	if githubMarkdownGenerator == nil {
		return nil
	}
	return githubMarkdownGenerator.WriteGithubJobSummary()
}

func NewGithubMarkdownGenerator(result *Result, title Operation) (markdownGenerator *MarkdownGenerator, cleanUp func() error, err error) {
	filename := os.Getenv(GithubEnvStepSummary)
	if filename == "" {
		wd, _ := os.Getwd()
		filename = path.Join(wd, "github-action-summary.md")
		// TODO change this to return nil, nil, nil
	}
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	cleanUp = func() error {
		return file.Close()
	}
	markdownGenerator = &MarkdownGenerator{file: file, result: result, operationTitle: title}
	return
}

func (m *MarkdownGenerator) WriteGithubJobSummary() (err error) {

	if m.result.SuccessCount() > 0 {
		var transferDetailsArray []*clientutils.FileTransferDetails
		for transferDetails := new(clientutils.FileTransferDetails); m.result.Reader().NextRecord(transferDetails) == nil; transferDetails = new(clientutils.FileTransferDetails) {
			transferDetailsArray = append(transferDetailsArray, transferDetails)
		}
		err = m.writeTable(m.operationTitle.String(), transferDetailsArray)
		if err != nil {
			return
		}
	}

	if m.result.FailCount() > 0 {
		a := "[🚨Error] Failed uploading %d artifacts.\n"
		err = m.writeHeader(fmt.Sprintf(a, m.result.FailCount()))
	}

	return
}

func (m *MarkdownGenerator) writeHeader(header string) error {
	_, err := m.file.WriteString(fmt.Sprintf("## %s\n", header))
	return err
}

func (m *MarkdownGenerator) writeSecondHeader(header string) error {
	_, err := m.file.WriteString(fmt.Sprintf("### %s\n", header))
	return err
}

func (m *MarkdownGenerator) writeString(content string) error {
	_, err := m.file.WriteString(content)
	return err
}

func (m *MarkdownGenerator) writeTable(header string, details []*clientutils.FileTransferDetails) error {
	// Write the table header
	_, err := m.file.WriteString(fmt.Sprintf("## %s\n", header))
	if err != nil {
		return err
	}

	// Write the table column names
	_, err = m.file.WriteString("| Source Path 📁   | Target Path 🎯  | Sha256 🔢  |\n| --- | --- |--- |\n")
	if err != nil {
		return err
	}

	// Write the table rows
	for _, detail := range details {
		line := fmt.Sprintf("| %s | %s | %s |\n", detail.SourcePath, path.Join(detail.RtUrl, detail.TargetPath), detail.Sha256)
		_, err = m.file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
