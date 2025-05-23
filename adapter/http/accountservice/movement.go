package accountservice

import (
	"encoding/json"
	"net/http"

	dto "github.com/pedro-scarelli/wheredidmymoneygo/core/dto/account/request"
)

func (service service) Movement(response http.ResponseWriter, request *http.Request) {
	movementRequestDTO, err := dto.FromJSONCreateMovementRequestDTO(request.Body)
	response.Header().Add("Content-Type", "application/json")

	if err != nil {
		response.WriteHeader(400)
		json.NewEncoder(response).Encode(map[string]string{"error": err.Error()})
		return
	}

	AccountID, err := GetAccountIDFromToken(request)
	if err != nil {
		response.WriteHeader(401)
		response.Write([]byte("Invalid token"))
		return
	}

	movementRequestDTO.AccountID = AccountID
	movementCreated, err := service.usecase.Movement(movementRequestDTO)
	if err != nil {
		response.WriteHeader(500)
		response.Write([]byte(err.Error()))
		return
	}

	response.WriteHeader(200)
	json.NewEncoder(response).Encode(movementCreated)
}
