package unit

import (
	"regexp"
	"testing"

	. "github.com/defenseunicorns/zarf/src/cmd"
)

func TestRemotePackage(t *testing.T) {
	t.Log("Unit: Remote Package")

	var (
		url = "https://github.com/defenseunicorns/zarf/tree/main/examples/config-file"
	)

	httpRegexTests(t, url)
	remotePackageTests(t, url)
}

func remotePackageTests(t *testing.T, url string) {
	CreatePackageCmd(nil, []string{url})
}

func httpRegexTests(t *testing.T, url string) {
	var (
		isRemoteUrlRegex  = regexp.MustCompile(UrlRegex)
		positiveTestCases = []string{
			url,
			"https://google.com",
			"http://www.google.com",
		}
		negativeTestCases = []string{
			"https://hkjjkhjkhjjk",
			"kfjklsdfjkdsfsdlkjfks",
			"23423423423",
			"ftp://github.com/defenseunicorns/zarf/tree/main/examples/config-file",
			"/Users/some/users/local/file/location",
			"/Users/some/users/local/file/location/file.zarf",
		}
	)

	for i := 0; i < len(positiveTestCases); i++ {
		if !isRemoteUrlRegex.MatchString(positiveTestCases[i]) {
			t.Fail()
		}
	}

	for i := 0; i < len(negativeTestCases); i++ {
		if isRemoteUrlRegex.MatchString(negativeTestCases[i]) {
			t.Fail()
		}
	}
}
