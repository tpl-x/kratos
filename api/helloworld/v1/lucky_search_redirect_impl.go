package v1

import "github.com/go-kratos/kratos/v2/transport/http"

var _ http.Redirector = (*LuckySearchResponse)(nil)

func (s *LuckySearchResponse) Redirect() (string, int) {
	return s.RedirectTo, int(s.StatusCode)
}
