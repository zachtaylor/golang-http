package hash

import "crypto/md5"

// Hash is a func(string) string
type Func = func(string) string

func NewMD5SaltHash(salt string) Func {
	return func(password string) string {
		sum := md5.Sum([]byte(password + salt))
		return string(sum[:])
	}
}
