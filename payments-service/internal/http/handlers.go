package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type callBackUriRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (s *HttpServer) handleCallBack(ctx *gin.Context) {
	var uri callBackUriRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "bad_request - missing transaction id"})
		return
	}

	_, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "invalid url"})
		return
	}

	var req any
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "no body received"})
		return
	}

	// map[
	// amount:50
	// forward_url:https://1336-105-163-156-6.ngrok-free.app/transaction/d8224676-dca9-4ded-99c8-cd3a9df70ec4
	// order_id:
	// phone_number:254718750145
	// remarks:Transaction processed successfully with reference: [SIA538W381]
	// result_code:0
	// third_party_trans_id:SIA538W381
	// transaction_date:November
	// transaction_reference:AIB081959911
	//]

	// transaction := Transaction{
	// 	Amount:            50,
	// 	ForwardURL:        "https://1336-105-163-156-6.ngrok-free.app/transaction/d8224676-dca9-4ded-99c8-cd3a9df70ec4",
	// 	OrderID:           "", // No order_id provided in the data
	// 	PhoneNumber:       "254718750145",
	// 	Remarks:           "Transaction processed successfully with reference: [SIA538W381]",
	// 	ResultCode:        0,
	// 	ThirdPartyTransID: "SIA538W381",
	// 	TransactionDate:   "November",
	// 	TransactionRef:    "AIB081959911",
	// }

	// this is for withdrawal
	// map[
	// amount:50
	// forward_url:https://1336-105-163-156-6.ngrok-free.app/transaction/215c54a9-1f80-4f08-ba68-b580cfb34003
	// order_id:
	// phone_number:2547*****145
	// remarks:Transaction processed successfully with reference: [SIA0446FKG]
	// result_code:0
	// third_party_trans_id:SIA0446FKG
	// transaction_date:2024-09-10 12:06:32.013711518 +0000 UTC m=+524359.621983232
	// transaction_reference:AIB120629278
	// ]

	// transaction := Transaction{
	// 	Amount:            50,
	// 	ForwardURL:        "https://1336-105-163-156-6.ngrok-free.app/transaction/215c54a9-1f80-4f08-ba68-b580cfb34003",
	// 	OrderID:           "", // No order_id provided in the data
	// 	PhoneNumber:       "2547*****145",
	// 	Remarks:           "Transaction processed successfully with reference: [SIA0446FKG]",
	// 	ResultCode:        0,
	// 	ThirdPartyTransID: "SIA0446FKG",
	// 	TransactionDate:   "2024-09-10 12:06:32.013711518 +0000 UTC m=+524359.621983232",
	// 	TransactionRef:    "AIB120629278",
	// }

	// probably change the transaction status after receiving the callback

	log.Println("request: ", req)
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
