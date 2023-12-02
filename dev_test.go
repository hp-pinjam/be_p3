package be_p3

import (
	"fmt"
	"testing"

	model "github.com/hp-pinjam/be_p3/model"
	modul "github.com/hp-pinjam/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = SetConnection("MONGOSTRING", "hppinjam")

// user
func TestRegister(t *testing.T) {
	var data model.User
	data.ID = primitive.NewObjectID()
	data.Email = "rijik@gmail.com"
	data.Username = "rijik"
	data.Role = "user"
	data.Password = "kepodah"

	err := modul.Register(mconn, "user", data)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success", data)
	}
}

// test login
func TestLogIn(t *testing.T) {
	var userdata model.User
	userdata.Username = "rijik"
	userdata.Password = "kepodah"
	user, status, err := modul.LogIn(mconn, "user", userdata)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	} else {
		fmt.Println("Login success", user)
	}
}

func TestUpdateUser(t *testing.T) {
	var data model.User
	data.Email = "rijik@gmail.com"
	data.Username = "rijik"
	data.Role = "admin"

	data.Password = "kepodah" // password tidak diubah

	id, err := primitive.ObjectIDFromHex("654a6513226d8ad245cd01ff")
	data.ID = id
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateUser(mconn, "user", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updateting document: %v", err)
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
	}
}

// test change password
func TestChangePassword(t *testing.T) {
	var data model.User
	data.Email = "rijik@gmail.com" // email tidak diubah
	data.Username = "rijik"        // username tidak diubah
	data.Role = "admin"            // role tidak diubah

	data.Password = "kepodah"

	// username := "dapskut123"

	_, status, err := modul.ChangePassword(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error updateting document: %v", err)
	} else {
		fmt.Println("Password berhasil diubah dengan username:", data.Username)
	}
}

// test delete user
func TestDeleteUser(t *testing.T) {
	username := "rijik"

	err := modul.DeleteUser(mconn, "user", username)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	} else {
		fmt.Println("Delete user success")
	}
}

func TestGetUserFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("6539d6c46700af5da789a678")
	anu, _ := modul.GetUserFromID(mconn, "user", id)
	fmt.Println(anu)
}

func TestGetUserFromUsername(t *testing.T) {
	anu, err := modul.GetUserFromUsername(mconn, "user", "rijik")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetUserFromEmail(t *testing.T) {
	anu, _ := modul.GetUserFromEmail(mconn, "user", "tejo@gmail.com")
	fmt.Println(anu)
}

func TestGetAllUser(t *testing.T) {
	anu, err := modul.GetAllUser(mconn, "user")
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
	fmt.Println(anu)
}

// hp
func TestInsertHp(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	var hpdata model.Hp
	hpdata.Title = "Iphone 15"
	hpdata.Description = "hp keluaran terbaru dari apple"
	hpdata.IsDone = true

	nama, err := modul.InsertHp(mconn, "hp", hpdata)
	if err != nil {
		t.Errorf("Error inserting hp: %v", err)
	}
	fmt.Println(nama)
}

func TestGetHpFromID(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "hppinjam")
	id, _ := primitive.ObjectIDFromHex("6548bce6c31c8ec3f02fa11d")
	anu := modul.GetHpFromID(mconn, "hp", id)
	fmt.Println(anu)
}

func TestGetHpList(t *testing.T) {
	anu, err := modul.GetHpList(mconn, "hp")
	if err != nil {
		t.Errorf("Error getting hp: %v", err)
		return
	}
	fmt.Println(anu)
}
