package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"vizyon-test/vizyon"
)

func main() {

	min := 100
	max := 300

	n := 100

	random := time.Now().UnixNano()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			amount := int64(rand.Intn(max-min+1) + min)
			orderId := strconv.FormatInt(random, 10)
			defer wg.Done()

			vizyon.PostAPIRequest(strconv.FormatInt(random, 10), vizyon.ProvisionRequest{
				ReturnUrl:   "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09-" + orderId,
				OrderId:     orderId,
				Amount:      amount,
				Installment: "1",
				Card: vizyon.ProvisionCardData{
					Cvc:         "000",
					ExpireMonth: "12",
					ExpireYear:  "2026",
					HolderName:  "TestKişisi",
					Number:      "4546711234567894",
				},
			}, true)
		}()
	}
	wg.Wait()

	fmt.Println("All processes completed!")
}
