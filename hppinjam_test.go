package be_p3

import (
	"fmt"
	"testing"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
)

// PASETO
func TestGeneratePrivateKeyPaseto(t *testing.T) {
	privateKey, publicKey := watoken.GenerateKey()
	fmt.Println(privateKey)
	fmt.Println(publicKey)
	hasil, err := watoken.Encode("wegotour", privateKey)
	fmt.Println(hasil, err)
}
func TestValidateToken(t *testing.T) {
	tokenstring := "eyJleHAiOiIyMDIzLTEwLTI2VDA1OjAyOjQ1WiIsImlhdCI6IjIwMjMtMTAtMjZUMDM6MDI6NDVaIiwiaWQiOiJkYWZmYSIsIm5iZiI6IjIwMjMtMTAtMjZUMDM6MDI6NDVaIn3cLq58WoqF4cfwdtKZiUas4-p4PVbwDaF4sa0QConAH_hZWT726D8" // Gantilah dengan token PASETO yang sesuai
	publicKey := "34984c89d5553bd07ced0b9ed6306cc010418a1758fae39e92bfce521ee7b44e"
	payload, _err := watoken.Decode(publicKey, tokenstring)
	if _err != nil {
		fmt.Println("expired token", _err)
	} else {
		fmt.Println("ID: ", payload.Id)
		fmt.Println("Di mulai: ", payload.Nbf)
		fmt.Println("Di buat: ", payload.Iat)
		fmt.Println("Expired: ", payload.Exp)
	}
}

// Hash Pass
func TestGeneratePasswordHash(t *testing.T) {
	password := "bisabis15"
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("Match:   ", match)
}

func TestHashFunction(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "wegotour")

	var userdata User
	userdata.Username = "dapskuy"
	userdata.Password = "kepoah"

	filter := bson.M{"username": userdata.Username}
	res := atdb.GetOneDoc[User](mconn, "user", filter)
	fmt.Println("Mongo User Result: ", res)
	hash, _ := HashPassword(userdata.Password)
	fmt.Println("Hash Password : ", hash)
	match := CheckPasswordHash(userdata.Password, res.Password)
	fmt.Println("Match:   ", match)
}

func TestInsertUser(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "wegotour")
	var userdata User
	userdata.Username = "daffa"
	userdata.Password = "kepoah"

	nama := InsertUser(mconn, "user", userdata)
	fmt.Println(nama)
}

func TestIsPasswordValid(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "wegotour")
	var userdata User
	userdata.Username = "dapskuy"
	userdata.Password = "kepoah"

	anu := IsPasswordValid(mconn, "user", userdata)
	fmt.Println(anu)
}
