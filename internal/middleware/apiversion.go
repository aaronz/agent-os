package middleware

import (
	"context"
	"net/http"
	"strings"
)

type APIVersion struct {
	version string
	major   int
	minor   int
}

func (v *APIVersion) String() string {
	return v.version
}

func (v *APIVersion) Major() int {
	return v.major
}

func (v *APIVersion) Minor() int {
	return v.minor
}

func (v *APIVersion) AtLeast(major, minor int) bool {
	if v.major > major {
		return true
	}
	if v.major == major && v.minor >= minor {
		return true
	}
	return false
}

func ParseAPIVersion(version string) *APIVersion {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")
	major := 1
	minor := 0

	if len(parts) > 0 && parts[0] != "" {
		if n, ok := parseInt(parts[0]); ok {
			major = n
		}
	}
	if len(parts) > 1 {
		if n, ok := parseInt(parts[1]); ok {
			minor = n
		}
	}

	return &APIVersion{
		version: version,
		major:   major,
		minor:   minor,
	}
}

func parseInt(s string) (int, bool) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

type apiVersionKey struct{}

func APIVersionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := "v1"

		if strings.HasPrefix(r.URL.Path, "/api/") {
			path := strings.TrimPrefix(r.URL.Path, "/api/")
			if idx := strings.Index(path, "/"); idx > 0 {
				version = path[:idx]
				r.URL.Path = "/api/" + path[idx+1:]
			}
		}

		ctx := context.WithValue(r.Context(), apiVersionKey{}, ParseAPIVersion(version))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAPIVersion(ctx context.Context) *APIVersion {
	if v := ctx.Value(apiVersionKey{}); v != nil {
		if version, ok := v.(*APIVersion); ok {
			return version
		}
	}
	return &APIVersion{version: "v1", major: 1, minor: 0}
}
