package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/EmilioCliff/payment-polling-app/payment-service/internal/services"
	"github.com/EmilioCliff/payment-polling-app/payment-service/pkg"
	"github.com/hibiken/asynq"
)

const SendWithdrawalRequestTask = "task:withdrawal_request"

func (distributor *RedisTaskDistributor) DistributeSendWithdrawalRequestTask(
	ctx context.Context,
	payload services.SendPaymentWithdrawalRequestPayload,
	opt ...asynq.Option,
) error {
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
	var taskPayload services.SendPaymentWithdrawalRequestPayload
	if err := json.Unmarshal(task.Payload(), &taskPayload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	payload := map[string]interface{}{
		"account_id":   taskPayload.PaydAccountID,
		"phone_number": taskPayload.PhoneNumber,
		"amount":       taskPayload.Amount,
		"narration":    taskPayload.Naration,
		"channel":      taskPayload.NetworkCode,
		"callback_url": fmt.Sprintf("%s/transaction/%v", processor.config.PAYD_CALLBACK_URL, taskPayload.TransactionID.String()),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	jsonString := strings.NewReader(string(jsonPayload))

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, processor.withdrawalUrl, jsonString)
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

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Failed to read response body: %w", err)
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(resBody, &responseData)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal response body: %w", err)
	}

	transactionReference, _ := responseData["correlator_id"].(string)
	message, _ := responseData["message"].(string)
	errorMessage, _ := responseData["error_message"].(string)

	if errorMessage != "" {
		message = "Payd Error: " + errorMessage
	}

	if res.StatusCode != http.StatusAccepted {
		processor.createTransaction(ctx, taskPayload, transactionReference, message, "failed")

		return fmt.Errorf("Request failed with status code: %d", res.StatusCode)
	}

	err = processor.createTransaction(ctx, taskPayload, transactionReference, message, "success")
	if err != nil {
		pkgErr, _ := err.(*pkg.Error)

		return fmt.Errorf("Failed to create transaction: %v\nWith error: %v", pkgErr.Message, pkgErr.Code)
	}

	return nil
}
