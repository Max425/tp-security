package handler

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"main/pkg/model/core"
	"main/pkg/model/dto"
	"net/http"
	"os/exec"
	"strings"
)

// @Summary get 1 request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} core.Request
// @Failure 500 {object} string
// @Router /requests/{uid} [get]
func (h *Handler) request(w http.ResponseWriter, r *http.Request) {
	UID, has := mux.Vars(r)["uid"]
	if !has {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	orders, err := h.repo.GetRequestByID(r.Context(), UID)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, orders)
}

// @Summary get all requests
// @Tags request
// @Accept  json
// @Produce  json
// @Success 200 {object} []core.Request
// @Failure 500 {object} string
// @Router /requests [get]
func (h *Handler) requests(w http.ResponseWriter, r *http.Request) {
	orders, err := h.repo.GetAllRequests(r.Context())
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, orders)
}

// @Summary repeat request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /repeat/{uid} [get]
func (h *Handler) repeat(w http.ResponseWriter, r *http.Request) {
	UID, has := mux.Vars(r)["uid"]
	if !has {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	orders, err := h.repo.GetRequestByID(r.Context(), UID)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	res, err := execute(orders)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, res)
}

// @Summary scan request
// @Tags request
// @Accept  json
// @Produce  json
// @Param uid path string true "request UID"
// @Success 200 {object} []string
// @Failure 500 {object} string
// @Router /scan/{uid} [get]
func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
	UID, has := mux.Vars(r)["uid"]
	if !has {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusBadRequest, "invalid params")
		return
	}
	orders, err := h.repo.GetRequestByID(r.Context(), UID)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	status, err := checkSQLInjection(orders)
	if err != nil {
		dto.NewErrorClientResponseDto(r.Context(), w, http.StatusInternalServerError, "unknown error")
		return
	}

	dto.NewSuccessClientResponseDto(r.Context(), w, status)
}

func checkSQLInjection(request core.Request) ([]string, error) {
	var res []string
	originalLength, err := execute(request)
	if err != nil {
		return res, err
	}

	parameters := make([]map[string]string, 0)
	parameters = append(parameters, request.GetParams)
	parameters = append(parameters, request.PostParams)
	parameters = append(parameters, request.Headers)
	parameters = append(parameters, request.Cookies)

	for _, params := range parameters {
		for key := range params {
			// Подстановка одинарной кавычки
			modifiedRequest := modifyRequest(request, key, "'")
			modifiedLength, _ := execute(modifiedRequest)
			if len(modifiedLength) != len(originalLength) {
				res = append(res, fmt.Sprintf("Параметр '%s' уязвим для SQL инъекций (одинарная кавычка)\n", key))
			}

			// Подстановка двойной кавычки
			modifiedRequest = modifyRequest(request, key, "\"")
			modifiedLength, _ = execute(modifiedRequest)
			if modifiedLength != originalLength {
				res = append(res, fmt.Sprintf("Параметр '%s' уязвим для SQL инъекций (двойная кавычка)\n", key))
			}
		}
	}
	if len(res) == 0 {
		return []string{"всё хорошо :)"}, nil
	} else {
		return res, nil
	}
}

func modifyRequest(request core.Request, key, value string) core.Request {
	modifiedRequest := request
	switch {
	case modifiedRequest.GetParams != nil:
		modifiedRequest.GetParams[key] += value
	case modifiedRequest.PostParams != nil:
		modifiedRequest.PostParams[key] += value
	case modifiedRequest.Headers != nil:
		modifiedRequest.Headers[key] += value
	case modifiedRequest.Cookies != nil:
		modifiedRequest.Cookies[key] += value
	}
	return modifiedRequest
}

func execute(request core.Request) (string, error) {
	var curlCommand bytes.Buffer

	// Формируем базовую часть команды curl
	curlCommand.WriteString("curl -x http://127.0.0.1:8080 ")
	if request.Method != "CONNECT" {
		curlCommand.WriteString("-X ")
		curlCommand.WriteString(request.Method)
	} else {
		request.Path = "https://" + strings.Split(request.Path, ":")[0]
	}
	for key, value := range request.GetParams {
		curlCommand.WriteString(fmt.Sprintf(" -G --data-urlencode \"%s=%s\"", key, value))
	}

	for key, value := range request.PostParams {
		curlCommand.WriteString(fmt.Sprintf(" -d \"%s=%s\"", key, value))
	}

	for key, value := range request.Headers {
		curlCommand.WriteString(fmt.Sprintf(" -H \"%s: %s\"", key, value))
	}

	for key, value := range request.Cookies {
		curlCommand.WriteString(fmt.Sprintf(" --cookie \"%s=%s\"", key, value))
	}

	curlCommand.WriteString(" " + request.Path)

	s := curlCommand.String()
	cmd := exec.Command("bash", "-c", s)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении команды curl: %v, вывод: %s", err, out)
	}

	res := strings.Split(string(out), "<html>")
	result := strings.Join(res[len(res)-1:], "<html>")
	return "<html>" + result, nil
}
