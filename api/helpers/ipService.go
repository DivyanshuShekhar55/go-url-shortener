package helper

import (
	"net/http"
)

func GetIPClient(w http.ResponseWriter, r *http.Request) (ip string) {

	ip = r.RemoteAddr

	// if behind a proxy or load balancer
	// check for X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ip = forwarded
	}

	return ip

}
