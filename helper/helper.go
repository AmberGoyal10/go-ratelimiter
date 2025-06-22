package helper

import (
	"errors"
	"net/http"
	"rate-limiter/constants"
)

func GetIpAddressFromHeaders(headers http.Header) (string, error) {
	var ip string
	for _, header := range constants.IP_HEADERS_PRIORITY_LIST {
		ip = headers.Get(header)
		if ip != "" {
			break
		}
	}

	if ip == "" {
		err := errors.New("ip address not found in headers")
		return "", err
	}

	return ip, nil
}
