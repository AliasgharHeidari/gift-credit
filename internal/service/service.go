package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/AliasgharHeidari/gift-credit/internal/model"
	"github.com/AliasgharHeidari/gift-credit/internal/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	InternalErr            = errors.New("internal server error")
	ErrGiftCodeUnavailable = errors.New("GiftCode unavailable")
	ErrGiftCodeOutOfUse    = errors.New("GiftCode is out of use")
	ErrGiftCodeAleadyUsed  = errors.New("GiftCode already used")
	ErrNotFound            = errors.New("GiftCode does not exist")
	ErrGiftCodeAleadyExist = errors.New("GiftCode already exist")
)

func UseGiftCode(req model.Input) (float64, error) {

	log.Print(req.Code)
	log.Print(req.Phone)

	var gift model.GiftCode
	DB := postgres.GetDB()

	tx := DB.Begin()

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("code = ?", req.Code).First(&gift).Error; err != nil {
		log.Println(err)
		tx.Rollback()
		return 0, err
	}

	if gift.IsActive == false {
		return 0, ErrGiftCodeUnavailable
	} else if gift.UsedCount >= 1000 {

		tx.Model(&model.GiftCode{}).Where("Code = ? ", req.Code).Update("is_active", false)
		tx.Rollback()
		return 0, ErrGiftCodeOutOfUse
	}
	var count int64

	err := tx.Model(&model.GiftCode{}).Where("Mobile_Number = ? AND Code = ?", req.Phone, req.Code).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return 0, InternalErr
	}

	if count > 0 {
		tx.Rollback()
		return 0, ErrGiftCodeAleadyUsed
	}

	// url topup request
	TopupReq()

	url := "http://localhost:9898/wallet/topup"

	body := map[string]interface{}{
		"mobile_number": req.Phone,
	}
	jsonBody, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to create request")
		tx.Rollback()
		return 0, InternalErr
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("failed to request wallet service, error:", err)
		tx.Rollback()
		return 0, InternalErr
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Print("wallet service returned status:", resp.StatusCode)
		tx.Rollback()
		return 0, InternalErr
	}

	tx.Model(&model.GiftCode{}).Where("code = ?", req.Code).Update("UsedCount", gorm.Expr("used_count + 1"))

	tx.Model(&model.GiftCode{}).Where("code = ?", req.Code).Update("Mobile_Number", req.Phone)

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

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return response.Balance, nil

}

func TopupReq()(){



	
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

func CreateGiftCode(NewGiftCode model.GiftCode) error {
	DB := postgres.GetDB()
	var count model.GiftCode
	result := DB.Model(&model.GiftCode{}).Where("Code = ?", NewGiftCode.Code).Find(&count)
	if result.RowsAffected != 0 {
		return ErrGiftCodeAleadyExist
	}
	if result.Error != nil {
		return InternalErr
	}

	DB.Save(&NewGiftCode)
	return nil

}
