package security // import "taylz.io/http/hash"

import "crypto/md5"

// F_Hash is a func(string) string
type F_Hash = func(string) string

// NewMD5SaltHash creates a Func type which salts input to crypto/md5.Sum
func NewMD5SaltHash(salt string) F_Hash {
	return func(password string) string {
		sum := md5.Sum([]byte(password + salt))
		return string(sum[:])
	}
}
