package main

import (
	"encoding/hex"
	"fmt"
	"github.com/miekg/pkcs11"
)

func main() {
	fmt.Println("Hello World.")

	p := pkcs11.New("/usr/lib/softhsm/libsofthsm2.so")
	if p == nil {
		fmt.Println("Error: p = nil.")
		return
	}
	if err := p.Initialize(); err != nil {
		panic(err)
	}

	defer p.Destroy()
	defer p.Finalize()

	tokenList, err := p.GetSlotList(false)
	if err != nil {
		panic(err)
	}

	fmt.Println("SlotNum: ", len(tokenList))

	s := tokenList[0] // slot id

	fmt.Println("Slot: 0 -> id: ", s)

	sh, err := p.OpenSession(s, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		panic(err)
	}

	fmt.Println("Session: ", sh)
	if err := p.Login(sh, pkcs11.CKU_USER, "1234"); err != nil {
		panic("user pin")
	}

	id, err := hex.DecodeString("20180706")
	if err != nil {
		panic(err)
	}
	if err := p.FindObjectsInit(sh, []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, "DMK"),
		pkcs11.NewAttribute(pkcs11.CKA_ID, id),
	}); err != nil {
		panic(err)
	}

	objs, _, err := p.FindObjects(sh, 1)
	if err != nil {
		panic(err)
	}
	p.FindObjectsFinal(sh)

	if len(objs) <= 0 {
		panic("object empty.")
	}

	aesKey := objs[0]

	/**
	var aesKey pkcs11.ObjectHandle
	for _, o := range objs {
		fmt.Println("Object: ", o)
		aesKey = o
		break
	}**/

	/**
	id, err := hex.DecodeString("9999")
	if err != nil {
		panic(err)
	}

	aesKey, err := p.GenerateKey(sh,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_AES_KEY_GEN, nil)},
		[]*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_SECRET_KEY),
			pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
			pkcs11.NewAttribute(pkcs11.CKA_SENSITIVE, true),
			pkcs11.NewAttribute(pkcs11.CKA_VALUE_LEN, 32),
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_AES),
			pkcs11.NewAttribute(pkcs11.CKA_LABEL, "AnhkTest"),
			pkcs11.NewAttribute(pkcs11.CKA_ID, id),
		})
	if err != nil {
		panic(err)
	}
	**/
	if err := p.EncryptInit(sh,
		[]*pkcs11.Mechanism{
			pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, []byte("1234567890123456")),
		}, aesKey); err != nil {
		panic(err)
	}

	e, err := p.Encrypt(sh, []byte("HelloWorld."))
	if err != nil {
		panic(err)
	}
	fmt.Println("encryptData: ", e)

	if err := p.DecryptInit(sh,
		[]*pkcs11.Mechanism{
			pkcs11.NewMechanism(pkcs11.CKM_AES_CBC_PAD, []byte("1234567890123456")),
		}, aesKey); err != nil {
		panic(err)
	}

	d, err := p.Decrypt(sh, e)
	if err != nil {
		panic(err)
	}

	fmt.Println("plainData: ", string(d))
	/*	if err := p.DestroyObject(sh, aesKey); err != nil {
			panic(err)
		}
	*/
}
