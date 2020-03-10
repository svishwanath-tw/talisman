package detector

import (
	"talisman/gitrepo"
	"talisman/talismanrc"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	customPatterns []talismanrc.PatternString
)

func TestShouldDetectPasswordPatterns(t *testing.T) {
	filename := "secret.txt"

	shouldPassDetectionOfSecretPattern(filename, []byte("\"password\" : UnsafePassword"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<password data=123> jdghfakjkdha</password>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<passphrase data=123> AasdfYlLKHKLasdKHAFKHSKmlahsdfLK</passphrase>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("<ConsumerKey>alksjdhfkjaklsdhflk12345adskjf</ConsumerKey>"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte("AWS key :"), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`BEGIN RSA PRIVATE KEY-----
	aghjdjadslgjagsfjlsgjalsgjaghjldasja
	-----END RSA PRIVATE KEY`), t)
	shouldPassDetectionOfSecretPattern(filename, []byte(`PWD=appropriate`), t)
}

func TestShouldIgnorePasswordPatterns(t *testing.T) {
	results := NewDetectionResults()
	content := []byte("\"password\" : UnsafePassword")
	filename := "secret.txt"
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	fileIgnoreConfig := talismanrc.FileIgnoreConfig{filename, "833b6c24c8c2c5c7e1663226dc401b29c005492dc76a1150fc0e0f07f29d4cc3", []string{"filecontent"}}
	ignores := &talismanrc.TalismanRC{FileIgnoreConfig: []talismanrc.FileIgnoreConfig{fileIgnoreConfig}}

	NewPatternDetector(customPatterns).Test(additions, ignores, results)
	assert.True(t, results.Successful(), "Expected file %s to be ignored by pattern", filename)
}

func shouldPassDetectionOfSecretPattern(filename string, content []byte, t *testing.T) {
	results := NewDetectionResults()
	additions := []gitrepo.Addition{gitrepo.NewAddition(filename, content)}
	NewPatternDetector(customPatterns).Test(additions, talismanRC, results)
	expected := "Potential secret pattern : " + string(content)
	assert.Equal(t, expected, getFailureMessage(results, additions))
	assert.Len(t, results.Results, 1)
}

func getFailureMessage(results *DetectionResults, additions []gitrepo.Addition) string {
	failureMessages := []string{}
	for _, failureDetails := range results.GetFailures(additions[0].Path) {
		failureMessages = append(failureMessages, failureDetails.Message)
	}
	if len(failureMessages) == 0 {
		return ""
	}
	return failureMessages[0]
}
