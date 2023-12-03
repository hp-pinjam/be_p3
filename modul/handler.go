package modul

import (
	"encoding/json"
	"net/http"
	"os"

	model "github.com/hp-pinjam/be_p3/model"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Responsed model.Credential
	Response  model.HpResponse
	datauser  model.User
	datahp    model.Hp
)

// user
func GCFHandlerGetAllUser(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	userlist, err := GetAllUser(mconn, collectionname)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Get User Success"
	Responsed.Data = userlist

	return GCFReturnStruct(Responsed)
}

func GCFHandlerGetUserByUsername(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	user, err := GetUserFromUsername(mconn, collectionname, username)
	if err != nil {
		Responsed.Message = "Error retrieving user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Hello user"
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerRegister(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	err = Register(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Register success"

	return GCFReturnStruct(Responsed)
}

func GCFHandlerLogIn(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := LogIn(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	tokenstring, err := watoken.Encode(user.UID, os.Getenv(PASETOPRIVATEKEYENV))
	if err != nil {
		Responsed.Message = "Gagal Encode Token :" + err.Error()

	} else {
		Responsed.Message = "Selamat Datang " + user.Username
		Responsed.Token = tokenstring
	}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerUpdateUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		Responsed.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Responsed.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := UpdateUser(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error updating user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Update success " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerChangePassword(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	user, _, err := ChangePassword(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error changing password: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Password change success for user " + user.Username
	Responsed.Data = []model.User{user}

	return GCFReturnStruct(Responsed)
}

func GCFHandlerDeleteUser(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Responsed.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Responsed.Message = "error parsing application/json1:"
		return GCFReturnStruct(Responsed)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Responsed.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Responsed)
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		Responsed.Message = "Missing 'username' parameter in the URL"
		return GCFReturnStruct(Responsed)
	}

	datauser.Username = username

	err = json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Responsed.Message = "error parsing application/json: " + err.Error()
	}

	_, err = DeleteUser(mconn, collectionname, datauser)
	if err != nil {
		Responsed.Message = "Error deleting user data: " + err.Error()
		return GCFReturnStruct(Responsed)
	}

	Responsed.Status = true
	Responsed.Message = "Delete user " + datauser.Username + " success"

	return GCFReturnStruct(Responsed)
}

// Hp
func GCFHandlerGetHp(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	// err = json.NewDecoder(r.Body).Decode(&datahp)
	// if err != nil {
	// 	Response.Message = "error parsing application/json3: " + err.Error()
	// 	return GCFReturnStruct(Response)
	// }

	hp, err := GetHpFromID(mconn, collectionname, ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Get hp success"
	Response.Data = []model.Hp{hp}

	return GCFReturnStruct(Response)
}

func GCFHandlerInsertHp(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	userInfo, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	err = json.NewDecoder(r.Body).Decode(&datahp)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	_, err = InsertHp(mconn, collectionname, datahp, userInfo.Id)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Insert hp success for " + datahp.Title
	Response.Data = []model.Hp{datahp}

	return GCFReturnStruct(Response)
}

func GCFHandlerUpdateHp(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	datahp.ID = ID

	err = json.NewDecoder(r.Body).Decode(&datahp)
	if err != nil {
		Response.Message = "error parsing application/json3: " + err.Error()
		return GCFReturnStruct(Response)
	}

	hp, _, err := UpdateHp(mconn, collectionname, datahp)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Update hp success"
	Response.Data = []model.Hp{hp}

	return GCFReturnStruct(Response)
}

func GCFHandlerDeleteHp(PASETOPUBLICKEY, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := MongoConnect(MONGOCONNSTRINGENV, dbname)
	Response.Status = false

	token := r.Header.Get("Authorization")
	if token == "" {
		Response.Message = "error parsing application/json1:"
		return GCFReturnStruct(Response)
	}

	_, err := watoken.Decode(os.Getenv(PASETOPUBLICKEY), token)
	if err != nil {
		Response.Message = "error parsing application/json2:" + err.Error() + ";" + token
		return GCFReturnStruct(Response)
	}

	id := r.URL.Query().Get("_id")
	if id == "" {
		Response.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Response.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(Response)
	}

	// err = json.NewDecoder(r.Body).Decode(&datahp)
	// if err != nil {
	// 	Response.Message = "error parsing application/json3: " + err.Error()
	// 	return GCFReturnStruct(Response)
	// }

	_, err = DeleteHp(mconn, collectionname, ID)
	if err != nil {
		Response.Message = err.Error()
		return GCFReturnStruct(Response)
	}

	Response.Status = true
	Response.Message = "Delete hp success"

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}
