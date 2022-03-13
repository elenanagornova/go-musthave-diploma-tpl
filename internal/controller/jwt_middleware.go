package controller

type CtxKey string

const UserCtxKey = CtxKey("UserID")

//func UserMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.URL.Path == "register" {
//			token, _ := auth.CreateToken()
//			http.SetCookie(w, &token)
//		}
//	})
//}
//
//func userUIDFromRequest(r *http.Request) string {
//	uid := r.Context().Value(UserCtxKey)
//	if uid == nil {
//		return ""
//	}
//	if userID, ok := uid.(string); ok {
//		return userID
//	}
//	return ""
//}
//
//func Welcome(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if strings.Contains(r.URL.Path, "register") {
//			next.ServeHTTP(w, r)
//		}
//	}}
