package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func maini() {

	router := gin.Default()

	router.POST("/webhook", func(ctx *gin.Context) {
		log.Printf("request body: %s", ctx.Request.Body)
	})

	router.GET("/send", func(ctx *gin.Context) {
		url := "https://api.mypayd.app/api/v1/payments"

		stringPayload := strings.NewReader(`{
			"first_name":     "emilio",
			"last_name":      "limo",
			"username":       "XI0LRwfl4uJIh0jBkbwe",
			"phone":          "254718750145",
			"amount":         5,
			"reason":         "testing",
			"email":          "emiliocliff@gmail.com",
			"location":       "Nairobi",
			"payment_method": "card",
			"callback_url":   "https://a1e0-41-90-186-109.ngrok-free.app/webhook",
		}`)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, stringPayload)

		req.Header.Set("Content-Type", "application/json")

		req.SetBasicAuth("CbCeJA4lYNxVMxb1tOC1", "YKIscOaKJxC4JB3qQ7V46VWe0gf39mHl7S66uBBK")
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	})

	router.GET("/payment-request", func(ctx *gin.Context) {
		url := "https://api.mypayd.app/api/v2/payments"

		stringPayload := strings.NewReader(`{
			"username": "emilio",
			"network_code": "63902",
			"amount": 10,
			"phone_number": "0718750145",
			"narration": "Payment for goods",
			"currency": "KES",
			"callback_url": "https://a1e0-41-90-186-109.ngrok-free.app/webhook"
		  }
		`)

		// data := map[string]interface{}{
		// 	"username":     "emiliocliff",
		// 	"network_code": "63902",
		// 	"amount":       10,
		// 	"phone_number": "0718750145",
		// 	"narration":    "Payment for goods",
		// 	"currency":     "KES",
		// 	"callback_url": "https://a1e0-41-90-186-109.ngrok-free.app/webhook",
		// }

		// jsonData, _ := json.Marshal(data)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, stringPayload)

		// req.Header.Set("Content-Type", "application/json")

		req.SetBasicAuth("CbCeJA4lYNxVMxb1tOC1", "YKIscOaKJxC4JB3qQ7V46VWe0gf39mHl7S66uBBK")
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	})

	router.GET("/card-details", func(ctx *gin.Context) {
		url := "https://api.mypayd.app/api/v2/payments"

		stringPayload := strings.NewReader(``)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, stringPayload)

		req.Header.Set("Content-Type", "application/json")

		req.SetBasicAuth("CbCeJA4lYNxVMxb1tOC1", "YKIscOaKJxC4JB3qQ7V46VWe0gf39mHl7S66uBBK")
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	})

	router.Run("0.0.0.0:3030")
}

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// )

// func main() {

// 	url := "https://api.mypayd.app/api/v1/payments"
// 	method := "POST"

// 	payload := strings.NewReader(`{
// "first_name" : "emilio",
// "last_name" : "limo",
// "amount": 10,
// "email": "emiliocliff@gmail.com",
// "location": "city",
// "username": "emilio",
// "payment_method": "card",
// "phone": "254718750145",
// "reason": "test",
// "callback_url": "http://your_callback_url.com"
// }`)

// 	client := &http.Client{}
// 	req, err := http.NewRequest(method, url, payload)

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	req.SetBasicAuth("CbCeJA4lYNxVMxb1tOC1", "YKIscOaKJxC4JB3qQ7V46VWe0gf39mHl7S66uBBK")
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer res.Body.Close()

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(string(body))
// }
