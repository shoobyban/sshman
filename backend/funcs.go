package backend

import (
	"crypto/sha1"
	"encoding/base64"
)

func checksum(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func match(aa, bb []string) bool {
	for _, a := range aa {
		for _, b := range bb {
			if a == b {
				return true
			}
		}
	}
	return false
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func updates(oldItems, newItems []string) (added []string, removed []string) {
	ma := make(map[string]struct{}, len(oldItems))
	mb := make(map[string]struct{}, len(newItems))
	for _, x := range newItems {
		mb[x] = struct{}{}
	}
	for _, x := range oldItems {
		if _, found := mb[x]; !found {
			removed = append(removed, x)
		}
		ma[x] = struct{}{}
	}
	for _, x := range newItems {
		if _, found := ma[x]; !found {
			added = append(added, x)
		}
	}
	return
}
