package handlers

import (
	"errors"
	"net/http"
)

type roomError struct {
	status  int
	message string
}

func (e roomError) Error() string { return e.message }

func writeRoomError(w http.ResponseWriter, err error) {
	status := http.StatusNotFound
	var e roomError
	if errors.As(err, &e) {
		status = e.status
	}
	http.Error(w, err.Error(), status)
}
