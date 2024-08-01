package app

import "encoding/json"

type initiatePaymentRequest struct {
	Amount        int64  `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}

func initiatePayment(data map[string]interface{}) []byte {
	// var payload initiatePaymentRequest
	// err := json.Unmarshal([]byte(data), &payload)
	// if err != nil {
	// 	return nil
	// }
	rsp, err := json.Marshal(data)
	if err != nil {
		return []byte("error marshalling request body")
	}

	return []byte(rsp)
}
