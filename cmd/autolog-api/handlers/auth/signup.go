package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"strings"

	"github.com/keola-dunn/autolog/internal/httputil"
	"github.com/keola-dunn/autolog/internal/logger"
	"github.com/keola-dunn/autolog/internal/service/user"
)

type signUpRequestBody struct {
	Username  string            `json:"username"`
	Email     string            `json:"email"`
	Password  string            `json:"password"`
	Name      string            `json:"name"`
	Questions []signupQuestions `json:"securityQuestions"`
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
	logEntry := logger.GetLogEntry(r)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logEntry.Error("failed to read sign up request body: %w", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if strings.TrimSpace(string(data)) == "" {
		httputil.RespondWithError(w, http.StatusBadRequest, "missing required signup request body")
		return
	}

	reqBody := signUpRequestBody{}

	if err := json.Unmarshal(data, &reqBody); err != nil {
		logEntry.Error("failed to unmarshal request body", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if valid, validationErr := reqBody.valid(); !valid {
		httputil.RespondWithError(w, http.StatusBadRequest, validationErr)
		return
	}
	ctx := r.Context()

	usernameExists, emailExists, err := a.userService.DoesUsernameOrEmailExist(ctx, reqBody.Username, reqBody.Email)
	if err != nil {
		logEntry.Error("failed to check if username or email exists", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	if usernameExists {
		httputil.RespondWithError(w, http.StatusConflict, "username already exists!")
		return
	}

	if emailExists {
		httputil.RespondWithError(w, http.StatusConflict, "email already exists!")
		return
	}

	var secQuestions = make([]user.UserSecurityQuestion, 0, len(reqBody.Questions))
	for _, q := range reqBody.Questions {
		secQuestions = append(secQuestions, user.UserSecurityQuestion{
			QuestionId: q.QuestionId,
			Answer:     q.Answer,
		})
	}

	userId, err := a.userService.CreateNewUser(ctx, user.CreateNewUserInput{
		Username:          reqBody.Username,
		Email:             reqBody.Email,
		Password:          reqBody.Password,
		SecurityQuestions: secQuestions,
		Role:              user.RoleUser,
	})
	if err != nil {
		logEntry.Error("failed to create new user", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	jwtToken, err := a.createJWT(userId)
	if err != nil {
		logEntry.Error("failed to create new user jwt", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	httputil.RespondWithJSON(w, http.StatusCreated, signUpResponse{
		JWT: jwtToken,
	})
}
