package aws_sdk

/* import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

//createHash File Crypto Encrypt & Decrypt
func createHash(byteStr []byte) []byte {

	hashVal := sha256.New()
	hashVal.Write(byteStr)
	bytes := hashVal.Sum(nil)
	return bytes
}

func encryptBytes(data []byte, passphrase []byte) ([]byte, error) {

	key := []byte(createHash(passphrase))
	key = key[0:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(string(data)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func decryptBytes(data []byte, passphrase []byte) ([]byte, error) {

	key := []byte(createHash(passphrase))
	key = key[0:16]

	text, _ := base64.StdEncoding.DecodeString(string(data))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, err
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	return text, nil
}
*/
