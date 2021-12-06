package response

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	Err() (err error)
	JSON(w http.ResponseWriter) (err error)
}

type responseImpl struct {
	err    error
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func Success(status string, data interface{}) (resp Response) {
	return &responseImpl{
		err:    nil,
		Status: status,
		Data:   data,
	}
}

func Error(status string, data interface{}, err error) (resp Response) {
	return &responseImpl{
		err:    err,
		Status: status,
		Data:   data,
	}
}

func (r *responseImpl) getStatusCode(status string) (statusCode int) {
	switch status {
	case StatusOK:
		return http.StatusOK
	case StatusCreated:
		return http.StatusCreated
	case StatusConflicted:
		return http.StatusConflict
	case StatusForbiddend:
		return http.StatusForbidden
	case StatusUnprocessabelEntity:
		return http.StatusUnprocessableEntity
	case StatusInvalidPayload:
		return http.StatusBadRequest
	case StatusUnexpectedError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func (r *responseImpl) Err() (err error) {
	return r.err
}

func (r *responseImpl) JSON(w http.ResponseWriter) (err error) {
	statusCode := r.getStatusCode(r.Status)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(r)
}
