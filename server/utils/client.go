package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
    // Check X-Forwarded-For header
    xff := r.Header.Get("X-Forwarded-For")
    if xff != "" {
        // X-Forwarded-For can contain multiple IPs, take the first one
        ips := strings.Split(xff, ",")
        if len(ips) > 0 {
            ip := strings.TrimSpace(ips[0])
            if net.ParseIP(ip) != nil {
                return ip
            }
        }
    }

    // Check X-Real-IP header
    xri := r.Header.Get("X-Real-IP")
    if xri != "" && net.ParseIP(xri) != nil {
        return xri
    }

    // Fall back to RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return r.RemoteAddr
    }
    
    return ip
}
