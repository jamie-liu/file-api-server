package utils

import (
    "crypto/rsa"
    "crypto/rand"
    "crypto/x509"
    "encoding/pem"
    "encoding/base64"
    "io/ioutil"
    "os"
)

func GenerateRSAKey(bits int, privateFile, publicFile string) error {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return err
    }

    X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
    privateFd, err := os.Create(privateFile)
    if err != nil {
        return err
    }
    defer privateFd.Close()
    privateBlock:= pem.Block{Type: "RSA Private Key", Bytes:X509PrivateKey}
    pem.Encode(privateFd, &privateBlock)

    publicKey := privateKey.PublicKey
    X509PublicKey,err := x509.MarshalPKIXPublicKey(&publicKey)
    if err != nil {
        return err
    }

    publicFd, err := os.Create(publicFile)
    if err != nil {
        return err
    }
    defer publicFd.Close()

    publicBlock := pem.Block{Type: "RSA Public Key", Bytes:X509PublicKey}
    return pem.Encode(publicFd, &publicBlock)
}

func RsaEncrypt(plainText, path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    block, _ := pem.Decode(content)
    if publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
        return "", err
    } else {
        publicKey,_ := publicKeyInterface.(*rsa.PublicKey)
        if buf, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(plainText)); err != nil {
            return "", err
        } else {
            return base64.StdEncoding.EncodeToString(buf), nil
        }
    }
}

func RsaDecrypt(cipherText, path string) (string, error) {
    content, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    block, _ := pem.Decode(content)
    decoded, err := base64.StdEncoding.DecodeString(cipherText)
    if err != nil {
        return "", err
    }

    if privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
        return "", err
    } else {
        buf, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decoded)
        return string(buf), err
    }
}
