package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"github.com/AliasgharHeidari/gift-credit/internal/repository/postgres"
)

var InternalErr = errors.New("internal server error")

func UseGiftCode(req model.Input) error {

	log.Print(req.Code)
	log.Print(req.Phone)

	var gift model.GiftCode
	DB := postgres.GetDB()

	if err := DB.Where("code = ?", req.Code).First(&gift).Error; err != nil {
		return err
	}

	url := "http://localhost:9898/wallet/gift"

	body := map[string]interface{}{
		"mobileNumber": req.Phone,
	}
	jsonBody, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to create request")
		return InternalErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("failed to request wallet service, error:", err)
		return InternalErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Print("wallet service returned status:", resp.StatusCode)
		return InternalErr
	}

	return nil

}


