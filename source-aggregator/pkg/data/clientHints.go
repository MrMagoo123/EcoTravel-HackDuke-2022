package data

import (
	"fmt"
	"strconv"
	"strings"
)

func parseSecChUa(userAgent string) string {
	if strings.Contains(userAgent, "Safari/605") || strings.Contains(userAgent, "Firefox") {
		return ""
	}

	version := "99"
	versionSplit := strings.Split(userAgent, "Chrome/")
	if len(versionSplit) > 1 {
		version = strings.Split(versionSplit[1], ".")[0]
	}

	matcher := ""
	if strings.Contains(userAgent, "Edg/") {
		matcher = "\"Chromium\";v=\"%s\", \" Not A;Brand\";v=\"99\", \"Microsoft Edge\";v=\"%s\""
	} else {
		versionInt, _ := strconv.Atoi(version)
		if versionInt == 106 {
			matcher = "\"Chromium\";v=\"%s\", \"Google Chrome\";v=\"%s\", \"Not;A=Brand\";v=\"99\""
		} else if versionInt == 105 {
			matcher = "\"Google Chrome\";v=\"%s\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"%s\""
		} else if versionInt > 102 {
			matcher = "\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"%s\", \"Chromium\";v=\"%s\""
		} else {
			matcher = "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"%s\", \"Google Chrome\";v=\"%s\""
		}
	}

	return fmt.Sprintf(matcher, version, version)
}

func parseSecChUaPlatform(userAgent string) string {
	if strings.Contains(userAgent, "Safari/605") || strings.Contains(userAgent, "Firefox") {
		return ""
	}

	if strings.Contains(userAgent, "Android") {
		return "\"Android\""
	} else if strings.Contains(userAgent, "Mac") {
		return "\"macOS\""
	} else if strings.Contains(userAgent, "Linux") {
		return "\"Linux\""
	} else {
		return "\"Windows\""
	}
}

func parseSecChUaMobile(userAgent string) string {
	if strings.Contains(userAgent, "Safari/605") || strings.Contains(userAgent, "Firefox") {
		return ""
	}

	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "Android") {
		return "?1"
	}

	return "?0"
}
