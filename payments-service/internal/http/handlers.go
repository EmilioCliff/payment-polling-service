package http

import (
	"net/http"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type callBackUriRequest struct {
	ID string `binding:"required" uri:"id"`
}

func (s *HttpServer) handleCallBack(ctx *gin.Context) {
	var uri callBackUriRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "bad_request - missing transaction id"})

		return
	}

	id, err := uuid.Parse(uri.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "invalid url"})

		return
	}

	var req any
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error_message": "no body received"})

		return
	}

	payload, _ := req.(map[string]interface{})

	resultCode, _ := payload["result_code"].(int)
	transactionRef, _ := payload["transaction_reference"].(string)
	remarks, _ := payload["remarks"].(string)

	status := false
	if resultCode == 0 {
		status = true
	}

	_, err = s.TransactionRepository.UpdateTransaction(ctx, id, repository.TransactionUpdate{
		Status:             status,
		PaydTransactionRef: transactionRef,
		Message:            remarks,
	})
	if err != nil {
		// log this data
		ctx.JSON(http.StatusInternalServerError, gin.H{"error_message": err.Error()})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// {
// 	"amount": 20,
// 	"forward_url": "https://603f-105-163-2-208.ngrok-free.app/callback",
// 	"order_id": "",
// 	"phone_number": "254718750145",
// 	"remarks": "Transaction processed successfully with reference: [SIR041D2V4]",
// 	"result_code": 0,
// 	"third_party_trans_id": "SIR041D2V4",
// 	"transaction_date": "November",
// 	"transaction_reference": "RIB181154135"
//   }
