package internal

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func (r *commandRouter) handlerB64decode(args string, user string, channel string) *string {
	result := new(string)
	b64decoded, err := base64.StdEncoding.DecodeString(args)
	if err != nil {
		*result = "Failed to decode string: " + err.Error()
		return result
	}
	*result = "Result: " + string(b64decoded)
	return result
}

func (r *commandRouter) handlerB64encode(args string, user string, channel string) *string {
	result := new(string)
	*result = "Result: " + base64.StdEncoding.EncodeToString([]byte(args))
	return result
}

func (r *commandRouter) handlerReverse(args string, user string, channel string) *string {
	runes := []rune(args)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	result := new(string)
	*result = string(runes)
	return result
}

func (r *commandRouter) handlerStrlen(args string, user string, channel string) *string {
	result := new(string)
	*result = fmt.Sprintf("Length: %d", len(args))
	return result
}

func (r *commandRouter) handlerMD5(args string, user string, channel string) *string {
	result := new(string)
	hasher := md5.New()
	hasher.Write([]byte(args))
	*result = "Result: " + hex.EncodeToString(hasher.Sum(nil))
	return result
}

func (r *commandRouter) handlerSHA1(args string, user string, channel string) *string {
	result := new(string)
	hasher := sha1.New()
	hasher.Write([]byte(args))
	*result = "Result: " + hex.EncodeToString(hasher.Sum(nil))
	return result
}

func (r *commandRouter) handlerSHA256(args string, user string, channel string) *string {
	result := new(string)
	hasher := sha256.New()
	hasher.Write([]byte(args))
	*result = "Result: " + hex.EncodeToString(hasher.Sum(nil))
	return result
}

func (r *commandRouter) handlerSHA512(args string, user string, channel string) *string {
	result := new(string)
	hasher := sha512.New()
	hasher.Write([]byte(args))
	*result = "Result: " + hex.EncodeToString(hasher.Sum(nil))
	return result
}
