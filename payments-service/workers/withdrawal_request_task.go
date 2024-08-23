package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
)

const SendWithdrawalRequestTask = "task:withdrawal_request"

func (distributor *RedisTaskDistributor) DistributeSendWithdrawalRequestTask(ctx context.Context, payload SendPaymentWithdrawalRequestPayload, opt ...asynq.Option) error {
	jsonWithdrawalRequestPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendPaymentRequestTask, jsonWithdrawalRequestPayload, opt...)
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

	log.Println("Processing task", task.Type(), taskPayload)

	url := "https://api.mypayd.app/api/v2/withdrawal"
	method := "POST"

	payload := strings.NewReader(`{
		"account_id": "your_account_id",
		"phone_number": "phone_number",
		"amount": 100,
		"narration": "Withdrawal request",
		"callback_url": "https://example.com/callback",
		"channel": "63902"
		}
	`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
