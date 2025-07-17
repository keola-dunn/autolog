package auth

import (
	"net/http"

	"github.com/keola-dunn/autolog/internal/httputil"
)

type getSecurityQuestionsResponse struct {
	Questions []SecurityQuestion `json:"questions"`
}

type SecurityQuestion struct {
	Question string `json:"question"`
	Id       string `json:"id"`
}

func (a *AuthHandler) GetSecurityQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := a.userService.GetSecurityQuestions(r.Context())
	if err != nil {
		a.logger.Error("failed to get security questions", err)
		httputil.RespondWithError(w, http.StatusInternalServerError, "")
		return
	}

	var resp getSecurityQuestionsResponse
	for _, q := range questions {
		resp.Questions = append(resp.Questions, SecurityQuestion{
			Question: q.Question,
			Id:       q.Id,
		})
	}

	httputil.RespondWithJSON(w, http.StatusOK, resp)
}
