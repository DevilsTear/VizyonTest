package vizyon

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
	"time"
)

type ProvisionRequest struct {
	ReturnUrl   string // Son kullanıcının ödemeyi onayladıktan sonra yönlendirileceği Url adresi
	OrderId     string // Bayinin oluşturduğu benzersiz sipariş numarasıdır.
	Amount      int64  // Karttan çekilecek tutardır. Kuruş kısmı . (nokta) ile yazılmalıdır. Örnek: 100.25
	Installment string // Taksit sayısıdır, tek çekim için 1 veya boş gönderilmelidir. Kabul edilen değerler: 1, 2, 3, 6, 9, 12
	Card        ProvisionCardData
}

type ProvisionResponse struct {
	Status       string // Yapılan API isteğinin başarılı olup olmadığını gösterir.
	ErrorCode    int64  // Yapılan istek başarısız olursa, hata kodu döner.
	ErrorType    string // Yapılan istek başarısız olursa, hata tipi döner.
	ErrorMessage string // Yapılan istek başarısız olursa, hata mesajı döner.
	ErrorDetail  string // Yapılan istek başarısız olursa, hata detayı döner.
	PaymentId    string // Ödeme Benzersiz Kimlik
	OrderId      string // İstek Yapılırken Gönderilen Sipariş Benzersiz Kimlik
	RedirectUrl  string // Ödeme Sayfasına Yönlendirme Url'i
	Language     string // İstek sonucunda dönen metinlerin dilini ayarlamak için kullanılır.
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PostAPIRequest(random string, provisionRequest ProvisionRequest, printResponse bool) (response ProvisionResponse, err error) {

	dataToEncrypt := fmt.Sprintf("%v%v%v", provisionRequest.OrderId, provisionRequest.Amount, provisionRequest.ReturnUrl)

	jsonBody, err := json.Marshal(provisionRequest)
	if err != nil {
		fmt.Printf("JSON marshal error: %v", err)
	}

	generatedHash, err := GenerateHash(random, dataToEncrypt)
	if err != nil {
		return ProvisionResponse{}, err
	}
	req, err := http.NewRequest("POST", vizyonProvisionURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return ProvisionResponse{}, err
	}
	req.Header.Set("ApiKey", apiKey)
	req.Header.Set("Random", random)
	req.Header.Set("Hash", generatedHash)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ProvisionResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if printResponse {
		fmt.Println("generatedHash:", generatedHash)
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return ProvisionResponse{}, err
	}

	return response, nil
}

func GenerateHash(random string, dataToEncrypt string) (string, error) {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(privateKey))

	// Write Data to it
	h.Write([]byte(dataToEncrypt))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	hash := fmt.Sprintf(`APIKey:%v&Random:%v&Signature:%v`, apiKey, random, sha)

	return b64.StdEncoding.EncodeToString([]byte(hash)), nil
}
