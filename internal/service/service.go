package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"github.com/AliasgharHeidari/gift-credit/internal/repository/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	InternalErr            = errors.New("internal server error")
	ErrGiftCodeUnavailable = errors.New("GiftCode unavailable")
	ErrGiftCodeOutOfUse    = errors.New("GiftCode is out of use")
	ErrGiftCodeAleadyUsed  = errors.New("GiftCode already used")
	ErrNotFound            = errors.New("GiftCode does not exist")
)

func UseGiftCode(req model.Input) (float64, error) {

	log.Print(req.Code)
	log.Print(req.Phone)

	var gift model.GiftCode
	DB := postgres.GetDB()

	if err := DB.Where("code = ?", req.Code).First(&gift).Error; err != nil {
		log.Println(err)
		return 0, err
	}

	if gift.IsActive == false {
		return 0, ErrGiftCodeUnavailable
	} else if gift.UsedCount >= 1000 {

		DB.Model(&model.GiftCode{}).Where("is_active = ?", gift.IsActive).Update("is_active", false)

		return 0, ErrGiftCodeOutOfUse
	}
	var count int64

	err := DB.Model(&model.GiftCodeUsage{}).Where("mobile_number = ?", req.Phone).Count(&count).Error
	if err != nil {
		return 0, InternalErr
	}

	if count > 0 {
		return 0, ErrGiftCodeAleadyUsed
	}

	url := "http://localhost:9898/wallet/topup"

	body := map[string]interface{}{
		"mobile_number": req.Phone,
	}
	jsonBody, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to create request")
		return 0, InternalErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("failed to request wallet service, error:", err)
		return 0, InternalErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Print("wallet service returned status:", resp.StatusCode)
		return 0, InternalErr
	}

	DB.Model(&model.GiftCode{}).Where("code = ?", req.Code).Update("UsedCount", gorm.Expr("used_count + 1"))

	NewRecord := model.GiftCodeUsage{
		MobileNumber: req.Phone,
	}
	DB.Create(&NewRecord)
	strPhone := strconv.Itoa(req.Phone)

	Url := "http://localhost:9898/wallet/" + strPhone

	resp, err = http.Get(Url)
	if err != nil {
		log.Println(err)
	}

	var response model.NewBalance
	Body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(Body))
	err = json.Unmarshal(Body, &response)

	log.Println(response.Balance)

	return response.Balance, nil

}

func GiftCodeStatus(GiftCode string) (model.GiftCode, error) {
	var GiftCodeStruct model.GiftCode

	DB := postgres.GetDB()

	result := DB.Model(&model.GiftCode{}).Where("Code = ? ", GiftCode).Find(&GiftCodeStruct)
	if result.Error != nil {
		return model.GiftCode{}, InternalErr
	}
	if result.RowsAffected == 0 {
		return model.GiftCode{}, ErrNotFound
	}

	return GiftCodeStruct, nil
}
