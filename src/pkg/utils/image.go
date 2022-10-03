package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// For further explanation see https://regex101.com/library/4Rl8mW and https://regex101.com/r/4Rl8mW/1
var hostParser = regexp.MustCompile(`(?im)([a-z0-9\-\_.]+)?(\/[a-z0-9\-.]+)?(:[\w\.\-\_]+)?$`)

// SwapHost Perform base url replacement and adds a sha1sum of the original url to the end of the src
func SwapHost(src string, targetHost string) (string, error) {
	targetImage, err := getTargetImageFromURL(src)
	return targetHost + "/" + targetImage, err
}

func getTargetImageFromURL(src string) (string, error) {
	submatches := hostParser.FindStringSubmatch(src)
	if len(submatches) == 0 {
		return "", fmt.Errorf("unable to get the targetImage from the provided source: %s", src)
	}

	// Combine (most) of the matches we obtained
	lastElementIndex := len(submatches) - 1
	targetImage := ""
	for _, match := range submatches[1:lastElementIndex] {
		targetImage += match
	}

	// Get a sha1sum of the src without a potential image tag
	tagMatcher := regexp.MustCompile(`(?im)(:[\w\.\-\_]+)?$`)
	srcWithoutTag := tagMatcher.ReplaceAllString(src, "")
	hasher := sha1.New()
	_, err := io.WriteString(hasher, srcWithoutTag)
	if err != nil {
		return "", fmt.Errorf("unable to get targetImage from the provided source: %w", err)
	}
	sha1Hash := hex.EncodeToString(hasher.Sum(nil))

	// Ensure we add the sha1sum before we apply an image tag
	if strings.HasPrefix(submatches[lastElementIndex], ":") {
		targetImage += "-" + sha1Hash + submatches[lastElementIndex]
	} else {
		targetImage += submatches[lastElementIndex] + "-" + sha1Hash
	}

	return targetImage, nil
}

// SwapHostWithoutSha Perform base url replacement but avoids adding a sha1sum of the original url.
func SwapHostWithoutSha(src string, targetHost string) (string, error) {
	submatches := hostParser.FindStringSubmatch(src)
	if len(submatches) == 0 {
		return "", fmt.Errorf("unable to get the targetImage from the provided source: %s", src)
	}

	return targetHost + "/" + submatches[0], nil
}