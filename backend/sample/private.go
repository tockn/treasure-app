package sample

import (
	"fmt"
	"log"
	"net/http"

	"github.com/voyagegroup/treasure-app/service"

	"github.com/voyagegroup/treasure-app/httputil"
)

type PrivateHandler struct {
	userService service.User
}

func NewPrivateHandler(us service.User) *PrivateHandler {
	return &PrivateHandler{
		userService: us,
	}
}

func (h *PrivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contextUser, err := httputil.GetUserFromContext(ctx)
	if err != nil {
		log.Print(err)
		WriteJSON(nil, w, http.StatusInternalServerError)
		return
	}
	user, err := h.userService.Get(ctx, contextUser.FirebaseUID)
	if err != nil {
		log.Printf("Show user failed: %s", err)
		WriteJSON(nil, w, http.StatusInternalServerError)
		return
	}
	resp := Response{
		Message: fmt.Sprintf("Hello %s from private endpoint! Your firebase uuid is %s", user.DisplayName, user.FirebaseUID),
	}
	WriteJSON(resp, w, http.StatusOK)
}
