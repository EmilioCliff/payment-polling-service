package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/postgres"
	"github.com/hibiken/asynq"
)

const SendWithdrawalRequestTask = "task:withdrawal_request"

func (distributor *RedisTaskDistributor) DistributeSendWithdrawalRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error {
	jsonWithdrawalRequestPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendWithdrawalRequestTask, jsonWithdrawalRequestPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: %s\n", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessWithdrawalRequestTask(ctx context.Context, task *asynq.Task) error {
	var taskPayload SendPaymentWithdrawalRequestPayload
	if err := json.Unmarshal(task.Payload(), &taskPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	url := "https://api.mypayd.app/api/v2/withdrawal"
	method := "POST"

	// enter manually the account id but will need this when registering the user
	payload := map[string]interface{}{
		"account_id":   "",
		"phone_number": taskPayload.PhoneNumber,
		"amount":       taskPayload.Amount,
		"narration":    taskPayload.Naration,
		"channel":      taskPayload.NetworkCode,
		"callback_url": fmt.Sprintf("https://1336-105-163-156-6.ngrok-free.app/transaction/%v", taskPayload.TransactionID.String()),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	jsonString := strings.NewReader(string(jsonPayload))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, jsonString)
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(taskPayload.PaydUsernameApiKey, taskPayload.PaydPasswordApiKey)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %w", err)
	}
	defer res.Body.Close()

	// check for other errors ie insufficient amount or low withdrawals < below 50
	// check if there was an error on the server side
	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf("Request failed with status code: %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %w", err)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(resBody, &responseData)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response body: %w", err)
	}

	// {
	// 	"status": 400,
	// 	"success": false,
	// 	"message": "mobile withdrawals have to be KES 50 and above"
	// }

	// get the withdrawal data to withdraw
	// create transaction but the action will be withdrawal
	transactionReference := responseData["message"].(string)

	// {
	// 	"success": true,
	// 	"correlator_id": "AIB053905490",
	// 	"message": "Transaction processed successfully.",
	// 	"status": "SUCCESS"
	// }

	_, pkgErr := processor.store.CreateTransactions(ctx, postgres.InitiatePaymentRequest{
		TransactionID:      taskPayload.TransactionID,
		PaydTransactionRef: transactionReference,
		UserID:             taskPayload.UserID,
		Action:             taskPayload.Action,
		Amount:             taskPayload.Amount,
		PhoneNumber:        taskPayload.PhoneNumber,
		NetworkCode:        taskPayload.NetworkCode,
		Naration:           taskPayload.Naration,
	})

	if pkgErr != nil {
		return fmt.Errorf("Failed to create transaction: %v\nWith error: %v", pkgErr.Message, pkgErr.Code)
	}

	return nil
}
