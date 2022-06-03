package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ProvisionRequest struct {
	ReturnUrl   string // Son kullanıcının ödemeyi onayladıktan sonra yönlendirileceği Url adresi
	OrderId     string // Bayinin oluşturduğu benzersiz sipariş numarasıdır.
	Amount      int64  // Karttan çekilecek tutardır. Kuruş kısmı . (nokta) ile yazılmalıdır. Örnek: 100.25
	Installment string // Taksit sayısıdır, tek çekim için 1 veya boş gönderilmelidir. Kabul edilen değerler: 1, 2, 3, 6, 9, 12
	Card        ProvisionCardData
}

type ProvisionCardData struct {
	HolderName  string // Kart sahibinin adı ve soyadıdır.
	Number      string // Kart numarasıdır.
	ExpireYear  string // Kartın son kullanma tarihi yılıdır.
	ExpireMonth string // Kartın son kullanma tarihi ayıdır.
	Cvc         string // Kartın arkasındaki güvenlik kodudur.
}

var vizyonProvisionURL string = "https://api.vizyonpos.com.tr/v1/payment/3d-secure/initialize"
var apiKey = "bc4fb4fe84714f82b829936eff596246"
var privateKey = "103f47fa9a5f4d64b80201491354880c"

func main() {
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 300

	random := time.Now().UnixNano()
	amount := int64(rand.Intn(max-min+1) + min)
	orderId := strconv.FormatInt(random, 10)

	PostAPIRequest(strconv.FormatInt(random, 10), ProvisionRequest{
		ReturnUrl:   "https://cc.mpay.software/cb.php?return_url=&p=31d565f77d64f8a453d6e598fb8c816026c9e72a3b0e98b3e7778bf17e805f09-" + orderId,
		OrderId:     orderId,
		Amount:      amount,
		Installment: "1",
		Card: ProvisionCardData{
			Cvc:         "123",
			ExpireMonth: "12",
			ExpireYear:  "2026",
			HolderName:  "Test Holder Name",
			Number:      "12345678901234",
		},
	})

}

func PostAPIRequest(random string, provisionRequest ProvisionRequest) {

	dataToEncrypt := fmt.Sprintf("%v%v%v", provisionRequest.OrderId, provisionRequest.Amount, provisionRequest.ReturnUrl)

	jsonBody, err := json.Marshal(provisionRequest)
	if err != nil {
		fmt.Printf("JSON marshal error: %v", err)
	}

	generatedHash, err := generateHash(random, dataToEncrypt)
	if err != nil {
		fmt.Println(err.Error())
	}
	req, err := http.NewRequest("POST", vizyonProvisionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("ApiKey", apiKey)
	req.Header.Set("Random", random)
	req.Header.Set("Hash", generatedHash)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func generateHash(random string, dataToEncrypt string) (string, error) {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(privateKey))

	// Write Data to it
	h.Write([]byte(dataToEncrypt))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	hash := fmt.Sprintf(`APIKey:%v&Random:%v&Signature:%v`, apiKey, random, sha)

	return b64.StdEncoding.EncodeToString([]byte(hash)), nil
}
