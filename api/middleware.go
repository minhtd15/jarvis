package api

import "net/http"

func AuthMiddleware(config Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(XApiKeyHeader)
			expectedApiKey := config.XApiKey

			if apiKey != expectedApiKey {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If xapikey is valid, continue
			next.ServeHTTP(w, r)
		})
	}
}
