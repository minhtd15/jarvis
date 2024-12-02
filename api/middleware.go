package api

import (
	"net/http"
)

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Thực hiện kiểm tra token ở đây nếu cần

			next.ServeHTTP(w, r) // Tiếp tục với handler tiếp theo
		})
		//return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	if strings.Contains(r.URL.Path, "/e/") || strings.Contains(r.URL.Path, "/web/") || strings.Contains(r.URL.Path, "/b/") {
		//		// url contains "/p/...", do not need token
		//		next.ServeHTTP(w, r)
		//		return
		//	}
		//
		//	authHeader := r.Header.Get("Authorization")
		//	if authHeader == "" {
		//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//		return
		//	}
		//
		//	tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		//	claims, err := jwtService.ValidateToken(tokenString)
		//	if err != nil {
		//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//		return
		//	}
		//
		//	// check whether the token is valid
		//	ctx := r.Context()
		//	ctx = context.WithValue(ctx, "username", claims.Username)
		//	ctx = context.WithValue(ctx, "user_id", claims.UserId)
		//	ctx = context.WithValue(ctx, "role", claims.Role)
		//	ctx = context.WithValue(ctx, "user_fullname", claims.UserFullName)
		//	//log.Infof("user name: %s \n userId: %s \n role: %s", claims.Username, claims.UserId, claims.Role)
		//	r = r.WithContext(ctx)
		//
		//	next.ServeHTTP(w, r)
		//})
	}
}
