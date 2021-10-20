package controllers

import (
	"net/http"

	"github.com/viadwimr/crowd_academy/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}