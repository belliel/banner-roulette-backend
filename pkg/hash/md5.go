package hash

import (
	"crypto/md5"
	"fmt"
	"time"
)

type MD5Hash struct {}

func NewMD5Hash() *MD5Hash {
	return &MD5Hash{}
}

func (M *MD5Hash) Hash() string {
	data := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	return fmt.Sprintf("%x", md5.Sum(data))
}

