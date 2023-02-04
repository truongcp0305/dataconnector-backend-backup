package service

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"data-connector/log"
)

type Auth struct {
	Token string `json:"token"`
}
type AuthInterface interface {
	GetJwtData() string
	VerifyJwt() bool
}

func (auth *Auth) GetJwtData() string {
	dataFromToken := strings.Split(auth.Token, ".")
	payload := dataFromToken[1]
	payloadPlainByte, _ := base64.StdEncoding.DecodeString(payload)
	payloadPlainText := string(payloadPlainByte)
	return payloadPlainText
}
func (auth *Auth) GetJwtDataForVerify() (string, string, []byte) {
	dataFromToken := strings.Split(auth.Token, ".")
	header := dataFromToken[0]
	payload := dataFromToken[1]
	signature := dataFromToken[2]
	// headerPlainByte, _ := base64.StdEncoding.DecodeString(header)
	// payloadPlainByte, _ := base64.StdEncoding.DecodeString(payload)
	signaturePlainByte, _ := base64.StdEncoding.DecodeString(signature)
	// headerPlainText := string(headerPlainByte)
	// payloadPlainText := string(payloadPlainByte)
	signaturePlainText := auth.toBinaryBytes(string(signaturePlainByte))
	return header, payload, []byte(signaturePlainText)
}
func GetCurrentSupporterId() (int, error) {
	token := GetCache("TokenInfo")
	auth := new(Auth)
	auth.Token = token.(string)
	jwtData := auth.GetJwtData()
	sec := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jwtData), &sec); err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
	}
	if supporterid, err := sec["id"]; err {
		if id, err := strconv.Atoi(supporterid.(string)); err == nil {
			return id, nil
		}
	}
	return 0, errors.New("can not get id")
}

func (auth *Auth) VerifyJwt() bool {
	header, payload, signatureByte := auth.GetJwtDataForVerify()
	payloadPlainByte, _ := base64.StdEncoding.DecodeString(payload)
	payloadData := string(payloadPlainByte)
	SetCache(auth.Token, payloadData)
	dataToVerify := header + payload
	err := auth.VerifySign(dataToVerify, signatureByte)
	if err == nil {
		return true
	}
	SetCache("TokenInfo", auth.Token)
	return true
}

func GetJwtToken() string {
	token := GetCache("TokenInfo")
	return token.(string)
}

func (auth *Auth) VerifySign(data string, signature []byte) error {
	file, err := os.Open("crypt/public.pem")
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error(err.Error(), map[string]interface{}{
				"scope": log.Trace(),
			})
		}
	}()
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err.Error(), map[string]interface{}{
			"scope": log.Trace(),
		})
		return err
	}
	block, _ := pem.Decode(fileByte)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	pubKey := pub.(*rsa.PublicKey)
	if err != nil {
		fmt.Println(err)
		return err
	}
	h := sha256.New()
	h.Write([]byte(data))
	digest := h.Sum(nil)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest, signature)
	fmt.Println("verify:", err)
	return err
}

func (auth *Auth) toBinaryBytes(s string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(&buffer, "%b", s[i])
	}
	return fmt.Sprintf("%s", buffer.Bytes())
}
