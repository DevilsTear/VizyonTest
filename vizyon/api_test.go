package vizyon

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"testing"
	"time"
	systemInfo "vizyon-test/system-info"
)

type testData struct {
	request ProvisionRequest
	random  string
}

var tests []testData

func init() {
	systemInfo.GetSystemInfo()
	min := 100
	max := 300

	n := 100

	for i := 0; i < n; i++ {
		random := time.Now().UnixNano()
		amount := int64(rand.Intn(max-min+1) + min)
		orderId := strconv.FormatInt(random, 10)
		tests = append(tests, testData{
			random: strconv.FormatInt(random, 10),
			request: ProvisionRequest{
				ReturnUrl:   "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09-" + orderId,
				OrderId:     orderId,
				Amount:      amount,
				Installment: "1",
				Card: ProvisionCardData{
					Cvc:         "000",
					ExpireMonth: "12",
					ExpireYear:  "2026",
					HolderName:  "TestKişisi",
					Number:      "4546711234567894",
				},
			},
		})
	}
}

func TestPostAPIRequestParallel(t *testing.T) {
	t.Parallel()
	for _, tc := range tests {
		tc := tc // capture range variable
		fmt.Printf("NumGoroutine is %d\n", runtime.NumGoroutine())
		response, err := PostAPIRequest(tc.random, tc.request, false)

		if err != nil {
			t.Error(err)
		}

		if response.ErrorCode != 0 && response.ErrorCode != 4011 {
			t.Error("Expected response.ErrorCode: 0, got ", response.ErrorCode)
		}
	}
}

func TestPostAPIRequest(t *testing.T) {
	response, err := PostAPIRequest("46699", ProvisionRequest{
		ReturnUrl:   "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09",
		OrderId:     "8051606",
		Amount:      200,
		Installment: "1",
		Card: ProvisionCardData{
			Cvc:         "000",
			ExpireMonth: "12",
			ExpireYear:  "2026",
			HolderName:  "TestKişisi",
			Number:      "4546711234567894",
		},
	}, false)

	if err != nil {
		t.Error(err)
	}

	if response.ErrorCode != 0 && response.ErrorCode != 4011 {
		t.Error("Expected response.ErrorCode: 0, got ", response.ErrorCode)
	}
}

func BenchmarkPostAPIRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PostAPIRequest("46699", ProvisionRequest{
			ReturnUrl:   "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09",
			OrderId:     "8051606",
			Amount:      200,
			Installment: "1",
			Card: ProvisionCardData{
				Cvc:         "000",
				ExpireMonth: "12",
				ExpireYear:  "2026",
				HolderName:  "TestKişisi",
				Number:      "4546711234567894",
			},
		}, false)
	}
}

func ExampleGenerateHash() {
	random, orderId, amount, returnUrl := "46699", "8051606", "200", "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09"
	dataToEncrypt := fmt.Sprintf("%v%v%v", orderId, amount, returnUrl)
	hash, _ := GenerateHash(random, dataToEncrypt)
	fmt.Println(hash)
	// Output:
	// QVBJS2V5OmJjNGZiNGZlODQ3MTRmODJiODI5OTM2ZWZmNTk2MjQ2JlJhbmRvbTo0NjY5OSZTaWduYXR1cmU6NTUyMDdiNzQwYWE5Nzk0NDU4OWQ1YzE3NjdhMWY5NjI0NDhmYWYxOTU3ODU5YTZkZjc4NjkxNTA0ZjljZDllYw==
}
