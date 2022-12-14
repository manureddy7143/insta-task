package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/manureddy7143/GolangStarter/source/service"
	"github.com/rs/zerolog/log"
)

const (
	err_code_a  = "err-400"
	err_code_b  = "err-401"
	err_code_c  = "err-404"
	http_code_a = http.StatusInternalServerError
	http_code_b = http.StatusBadRequest
	http_code_c = http.StatusNotFound
	http_code_d = http.StatusUnauthorized
	http_code_e = http.StatusConflict
)

func errorHandling(c *gin.Context, msg string, err error) {
	response := func(errMsg string, errCode string, http_code int) {
		log.Error().Msgf("%s: %s", errMsg, err.Error())
		c.JSON(http_code, dto.ErrorDTO{ErrorCode: errCode, ErrorMessage: errMsg})
	}

	switch msg {
	case "err-5000":
		response("Internal Server Error", "err-500", http_code_a)
		return
	case "err-5001":
		response("Binding Error", "err-500", http_code_a)
		return
	case "err-4000":
		response("Bad Request", err_code_a, http_code_b)
		return
	case "err-4001":
		response("User dont exists", err_code_a, http_code_b)
		return
	case "err-4002":
		response("Invalid Credentials", err_code_b, http_code_d)
		return

	case "err-4003":
		response("Invalid Password String", err_code_a, http_code_b)
		return
	default:
		log.Error().Msgf("Unknown_Error: %s", err.Error())
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorDTO{ErrorCode: "err-422", ErrorMessage: "Unknown Error"})
		return
	}
}
