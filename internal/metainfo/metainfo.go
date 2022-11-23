package metainfo

import (
	"fmt"
	"strings"
)

func strOrNA(s string) string {
	if strings.TrimSpace(s) == "" {
		return "N/A"
	}
	return s
}

func PrintBuildInfo(buildVersion string, buildDate string, buildCommit string) {
	fmt.Printf("buildVersion=%s buildDate=%s buildCommit=%s\n", strOrNA(buildVersion), strOrNA(buildDate), strOrNA(buildCommit))
}
