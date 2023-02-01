package lib

import (
	"crypto/x509"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
)

// HashCalculation calculat hash
func HashCalculation(h hash.Hash, val string) string {
	h.Write([]byte(val))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ReadCertPool get CertPool by crt file
func ReadCertPool(crt string) *x509.CertPool {
	_crt, err := ioutil.ReadFile(crt)
	if err != nil {
		log.Fatalln("Read crt file failed:", err.Error())
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(_crt) {
		log.Fatalln("Load crt file failed.")
	}
	return pool
}
