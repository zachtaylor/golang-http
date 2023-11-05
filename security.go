package http // import "taylz.io/http/security"

import "crypto/md5"

// NewMD5SaltHashFunc creates a Func type which salts input to crypto/md5.Sum
func NewMD5SaltHashFunc(salt string) func(string) string {
	return func(password string) string {
		sum := md5.Sum([]byte(password + salt))
		return string(sum[:])
	}
}
