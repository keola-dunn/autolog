package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type signUpRequestBody struct {
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Password  string            `json:"password"`
	Name      string            `json:"name"`
	Questions []signupQuestions `json:"questions"`
}

type signupQuestions struct {
	QuestionId string `json:"questionId"`
	Answer     string `json:"answer"`
}

func (s *signUpRequestBody) valid() (bool, string) {
	if _, err := mail.ParseAddress(s.Email); err != nil {
		return false, "email address is invalid"
	}
	if strings.TrimSpace(s.Username) == "" || len(s.Username) > 64 {
		return false, "username is missing or invalid"
	}
	if strings.TrimSpace(s.Password) == "" || len(s.Password) < 3 {
		return false, "missing or too short password"
	}
	if len(s.Questions) < 3 {
		return false, "missing security questions"
	}
	return true, ""
}

type signUpResponse struct {
	JWT string `json:"jwt"`
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		a.logger.Error("failed to read sign up request body: %w", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if strings.TrimSpace(string(data)) == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "missing required signup request body")
		return
	}

	reqBody := signUpRequestBody{}

	if err := json.Unmarshal(data, &reqBody); err != nil {
		a.logger.Error("failed to unmarshal request body", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if valid, validationErr := reqBody.valid(); !valid {
		httputil.RespondWithError(w, http.StatusBadRequest, validationErr)
		return
	}

	userId, err := a.userService.CreateNewUser(r.Context(), user.CreateNewUserInput{
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		// TODO: handle if username or email already exists
		a.logger.Error("failed to create new user", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	jwtToken, err := a.createJWT(userId)
	if err != nil {
		a.logger.Error("failed to create new user jwt", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, signUpResponse{
		JWT: jwtToken,
	})
}
