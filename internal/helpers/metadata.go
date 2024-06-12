package helpers

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// Do a case-insensitive key search within the given metadata
func MetadataGetCi(md metadata.MD, k string) []string {
	k = strings.ToLower(k);
	return md.Get(k)
}