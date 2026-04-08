package lib

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"hash"
	"log"
	"net/http"
	"os"
)

// HashCalculation calculat hash
func HashCalculation(h hash.Hash, val string) string {
	h.Write([]byte(val))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ReadCertPool get CertPool by crt file
func ReadCertPool(crt string) *x509.CertPool {
	_crt, err := os.ReadFile(crt)
	if err != nil {
		log.Fatalln("Read crt file failed:", err.Error())
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(_crt) {
		log.Fatalln("Load crt file failed.")
	}
	return pool
}

// HttpWriteJSON write JSON response
func HttpWriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if statusCode != 0 {
		w.WriteHeader(statusCode)
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
}
