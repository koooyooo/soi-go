package commons

import "strings"

func DefaultName(uri string) string {
	idxURIServerStart := strings.Index(uri, "://") + len("://")
	fromStartUri := uri[idxURIServerStart:]
	idxFromStartServerEnd := strings.Index(fromStartUri, "/")
	if idxFromStartServerEnd == -1 {
		idxFromStartServerEnd = len(fromStartUri)
	}
	server := fromStartUri[:idxFromStartServerEnd]
	wwwLessServer := strings.TrimPrefix(server, "www.")
	return wwwLessServer
}

func TrimElements(before []string) []string {
	var after []string
	for _, b := range before {
		after = append(after, strings.TrimSpace(b))
	}
	return after
}
