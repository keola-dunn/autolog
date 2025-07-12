package signup

import "net/http"

type SignupHandler struct{}

func NewSignupHandler() *SignupHandler {
	return &SignupHandler{}
}

func (h *SignupHandler) SignUp(w http.ResponseWriter, r *http.Request) {

}
