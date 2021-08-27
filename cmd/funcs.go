package cmd

import (
	"crypto/sha1"
	"encoding/base64"
)

func checksum(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func findByEmail(conf *config, email string) *user {
	for _, u := range conf.Users {
		if u.Email == email {
			return &u
		}
	}
	return nil
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
