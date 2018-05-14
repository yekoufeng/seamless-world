package login

import (
	"net/http"
	"regexp"

	"github.com/ant0ine/go-json-rest/rest"
)

// RequireHostMiddleware 只允许通过域名访问的请求
type RequireHostMiddleware struct{}

// MiddlewareFunc 具体实现
func (mw *RequireHostMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {
		re, _ := regexp.Compile(`([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`)
		if re.MatchString(r.Host) {
			rest.Error(w, "不允许直接访问", http.StatusForbidden)
			return
		}

		handler(w, r)
	}
}
