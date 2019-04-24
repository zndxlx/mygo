package service

import (
    "common/crypto/padding"
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    log "github.com/thinkboy/log4go"
    "math/rand"
    "time"
)

func AesCBCEncrypt(input []byte, key string) (output string, err error) {
    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        log.Error("aes.NewCipher failed, err=%s", err.Error())
        return
    }

    blockSize := block.BlockSize()
    padInput := padding.PKCS5.Padding(input, blockSize)

    aesOutput := make([]byte, len(padInput)+16)

    aesIv := aesOutput[:16]
    crypted := aesOutput[16:]

    randBytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < 16; i++ {
        aesIv[i] = randBytes[r.Intn(len(randBytes))]
    }

    blockMode := cipher.NewCBCEncrypter(block, aesIv)
    blockMode.CryptBlocks(crypted, padInput)

    output = base64.StdEncoding.EncodeToString(aesOutput)

    log.Info("input=%s, key=%s, aesIv=%s, output=%s", string(input), key,
        string(aesIv), output)
    return
}

func AesCBCDecrypt(input string, key string) (output []byte, err error) {
    decodeBytes, err := base64.StdEncoding.DecodeString(input)
    if err != nil {
        log.Error("base64 decode failed, err=%s", err.Error())
        return
    }

    block, err := aes.NewCipher([]byte(key))
    if err != nil {
        log.Error("aes.NewCipher failed, err=%s", err.Error())
        return
    }

    blockSize := block.BlockSize()
    blockMode := cipher.NewCBCDecrypter(block, decodeBytes[:16])
    output = make([]byte, len(decodeBytes)-16)

    blockMode.CryptBlocks(output, decodeBytes[16:])
    output, err = padding.PKCS5.Unpadding(output, blockSize)

    log.Info("input=%s, key=%s, output=%s", input, key, string(output))

    return
}
