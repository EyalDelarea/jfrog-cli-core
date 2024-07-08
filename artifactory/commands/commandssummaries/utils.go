package commandssummaries

import "strings"

// Protocols are being replaced due masking issue of secrets which are not secrets
// This function should only be used inside the context of command summaries.
// As the links are JFrog platform links, we can safely replace https with http.
func replaceProtocol(input string) string {
	if strings.HasPrefix(input, "https://") {
		return "http://" + input[len("https://"):]
	}
	return input
}
