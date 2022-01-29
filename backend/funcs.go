package backend

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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

func match(group1, group2 []string) bool {
	for _, a := range group1 {
		for _, b := range group2 {
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

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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

// Read and split SSH Public Key file into []string parts
// order will be key type, key, comment
// we assume comment is a user name or email
func readKeyFile(filename string) ([]string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error: error reading public key file: '%s' %v", filename, err)
	}
	return SplitParts(string(b))
}

// SplitParts splits a public key into []string parts
// order will be key type, key, comment
// we assume comment is a user name or email
func SplitParts(content string) ([]string, error) {
	parts := strings.Split(strings.TrimSuffix(content, "\n"), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("error: not a proper public key file")
	}
	return parts, nil
}

func JSON(data interface{}) string {
	// marshal data into json
	bs, _ := json.MarshalIndent(data, "", "  ")
	return string(bs)
}
