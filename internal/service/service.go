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
	// checks if GiftCode is avalible
	err := CheckIsActive(&gift)
	if err != nil {
		tx.Rollback()
		return 0, ErrGiftCodeOutOfUse
	}
	// checks if input number already used the giftcode
	res, err := CheckAlreadyUsed(tx, req)
	if err != nil {
		tx.Rollback()
		return 0, InternalErr
	}
	if res == false {
		tx.Rollback()
		return 0, ErrGiftCodeAleadyUsed
	}
	// requests wallet service to increase balance
	resp, err := TopUpRequest(req)
	if err != nil {
		tx.Rollback()
		return 0, InternalErr
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		log.Print("wallet service returned status:", resp.StatusCode)
		tx.Rollback()
		return 0, InternalErr
	}
	// Updates fields of usages for limits
	UpdateUsages(tx, req)

	// request wallet service for the new balance (after applying GiftCode)

	response := GetNewBalance(req)

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return response.Balance, nil

}

// following funcs are related to func UseGiftCode

func CheckIsActive(gift *model.GiftCode) error {
	var err = errors.New("err")
	if gift.IsActive == false {
		return err
	} else if gift.UsedCount >= 1000 {
		return err
	}
	return nil
}

func CheckAlreadyUsed(tx *gorm.DB, req model.Input) (bool, error) {
	var count int64
	err := tx.Model(&model.GiftCode{}).Where("Code = ? AND Mobile_Number = ?", req.Code, req.Phone).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}
	return true, nil
}

func TopUpRequest(req model.Input) (*http.Response, error) {

	url := "http://localhost:9898/wallet/topup"

	body := map[string]interface{}{
		"mobile_number": req.Phone,
	}
	jsonBody, _ := json.Marshal(body)

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Print("failed to create request")
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("failed to request wallet service, error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil

}

func UpdateUsages(tx *gorm.DB, req model.Input) {
	tx.Model(&model.GiftCode{}).Where("code = ?", req.Code).Update("UsedCount", gorm.Expr("used_count + 1"))
	tx.Model(&model.GiftCode{}).Where("code = ?", req.Code).Update("Mobile_Number", req.Phone)
}

func GetNewBalance(req model.Input)(model.NewBalance) {

	strPhone := strconv.Itoa(req.Phone)

	Url := "http://localhost:9898/wallet/" + strPhone

	resp, err := http.Get(Url)
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
	return response
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
