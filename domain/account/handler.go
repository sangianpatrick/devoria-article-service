package account

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sangianpatrick/devoria-article-service/middleware"
	"github.com/sangianpatrick/devoria-article-service/response"
)

type AccountHTTPHandler struct {
	Validate *validator.Validate
	Usecase  AccountUsecase
}

func NewAccountHTTPHandler(
	router *mux.Router,
	basicAuthMiddleware middleware.RouteMiddleware,
	validate *validator.Validate,
	usecase AccountUsecase,
) {
	handler := &AccountHTTPHandler{
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/accounts/registration", basicAuthMiddleware.Verify(handler.Register)).Methods(http.MethodPost)

}

func (handler *AccountHTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var resp response.Response
	var params AccountRegistrationRequest
	var ctx = r.Context()

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		resp = response.Error(response.StatusUnprocessabelEntity, nil, err)
		resp.JSON(w)
		return
	}

	err = handler.Validate.StructCtx(ctx, params)
	if err != nil {
		resp = response.Error(response.StatusInvalidPayload, nil, err)
		resp.JSON(w)
		return
	}

	resp = handler.Usecase.Register(ctx, params)
	resp.JSON(w)
}
