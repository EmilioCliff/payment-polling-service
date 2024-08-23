package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const SendPaymentRequestTask = "task:payment_request"

type SendPaymentWithdrawalRequestPayload struct {
	TransactionID      uuid.UUID `json:"transaction_id"`
	UserID             int64     `json:"user_id"`
	Action             string    `json:"action"`
	Amount             int64     `json:"amount"`
	PhoneNumber        string    `json:"phone_number"`
	NetworkCode        string    `json:"network_code"`
	Naration           string    `json:"naration"`
	PaydUsername       string    `json:"payd_username"`
	PaydPasswordApiKey string    `json:"payd_password_api_key"`
	PaydUsernameApiKey string    `json:"payd_username_api_key"`
}

func (distributor *RedisTaskDistributor) DistributeSendPaymentRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error {
	jsonPaymentRequestPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendPaymentRequestTask, jsonPaymentRequestPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Enqueued task: %s\n", info.ID)
	return nil
}

func (processor *RedisTaskProcessor) ProcessPaymentRequestTask(ctx context.Context, task *asynq.Task) error {
	var taskPayload SendPaymentWithdrawalRequestPayload
	if err := json.Unmarshal(task.Payload(), &taskPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	log.Println("Processing task: ", task.Type(), taskPayload)

	url := "https://api.mypayd.app/api/v2/payments"
	method := "POST"

	payload := map[string]interface{}{
		"username":     taskPayload.PaydUsername,
		"network_code": taskPayload.NetworkCode,
		"amount":       taskPayload.Amount,
		"phone_number": taskPayload.PhoneNumber,
		"narration":    taskPayload.Naration,
		"currency":     "KES",
		"callback_url": "https://4343-41-90-180-76.ngrok-free.app/" + taskPayload.TransactionID.String(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonPayload))

	req.SetBasicAuth(taskPayload.PaydUsernameApiKey, taskPayload.PaydPasswordApiKey)

	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 202 {
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

	transactionReference := responseData["merchantRequestID"].(string)

	createTransaction := createTransactionData{
		TransactionID:      taskPayload.TransactionID,
		PaydTransactionRef: transactionReference,
		UserID:             taskPayload.UserID,
		Action:             taskPayload.Action,
		Amount:             taskPayload.Amount,
		PhoneNumber:        taskPayload.PhoneNumber,
		NetworkCode:        taskPayload.NetworkCode,
		Naration:           taskPayload.Naration,
	}

	err = processor.createTransaction(ctx, createTransaction)
	if err != nil {
		return fmt.Errorf("Failed to create transaction: %w", err)
	}

	return nil
}
