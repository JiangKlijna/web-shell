package lib

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// seed Secret.seed
var seed = strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.FormatUint(rand.New(rand.NewSource(time.Now().UnixNano())).Uint64(), 36)

// GenerateSecret Get secret
// secret = sha224(seed+Version+GOOS+GOARCH+date+clientIP+userAgent).reverse()
func GenerateSecret(clientIP, userAgent string) string {
	date := time.Now().Format("2006-01-02")
	return ReverseString(HashCalculation(sha256.New224(), seed+runtime.Version()+runtime.GOOS+runtime.GOARCH+date+clientIP+userAgent))
}

// GenerateToken Get token
// token = md5(secret+md5(username+secret+password)+secret)
func GenerateToken(username, password, secret string) string {
	return HashCalculation(md5.New(), secret+HashCalculation(md5.New(), username+secret+password)+secret)
}

// GeneratePath Get path
// path = sha512(secret.reverse()^5+token.reverse()^5).reverse()
func GeneratePath(secret, token string) string {
	return ReverseString(HashCalculation(sha512.New(), strings.Repeat(ReverseString(secret), 5)+strings.Repeat(ReverseString(token), 5)))
}

// GenerateAll Get secret token path
func GenerateAll(username, password, clientIP, userAgent string) (string, string, string) {
	secret := GenerateSecret(clientIP, userAgent)
	token := GenerateToken(username, password, secret)
	path := GeneratePath(secret, token)
	return secret, token, path
}
