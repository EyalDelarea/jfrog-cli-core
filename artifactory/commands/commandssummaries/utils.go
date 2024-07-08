package commandssummaries

import (
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strings"
)

// Protocols are being replaced due masking issue of secrets which are not secrets
// This function should only be used inside the context of command summaries.
// As the links are JFrog platform links, we can safely replace https with http.
func replaceProtocol(input string) string {
	if strings.HasPrefix(input, "https://") {
		log.Info("Replacing https with http in the URL")
		log.Info("Input URL: " + input)
		log.Info("Output URL: " + "http://" + input[len("https://"):])
		return "http://" + input[len("https://"):]
	}
	log.Info("Input is not a URL, skipping protocol replacement")
	return input
}
