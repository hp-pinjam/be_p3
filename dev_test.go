package be_p3

import (
	"fmt"
	"testing"

	model "github.com/hp-pinjam/be_p3/model"
	modul "github.com/hp-pinjam/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = modul.MongoConnect("MONGOSTRING", "hppinjam")

// user
func TestRegister(t *testing.T) {
	var data model.User
	data.Email = "rijik@gmail.com"
	data.Username = "rijik"
	// data.Role = "user"
	data.Password = "secret"
	data.ConfirmPassword = "secret"

	err := modul.Register(mconn, "user", data)
	if err != nil {
		t.Errorf("Error registering user: %v", err)
	} else {
		fmt.Println("Register success", data.Username)
	}
}

// test login
func TestLogIn(t *testing.T) {
	var data model.User
	data.Username = "rijik"
	data.Password = "secret"
	data.Role = "user"

	user, status, err := modul.LogIn(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error logging in user: %v", err)
	} else {
		fmt.Println("Login success", user)
	}
}

func TestUpdateUser(t *testing.T) {
	var data model.User
	data.Email = "rijikhaura@gmail.com"
	data.Username = "rijik"

	id := "656c3f638442be4a7c185a09"
	ID, err := primitive.ObjectIDFromHex(id)
	data.ID = ID
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
	data.Password = "secrets"
	data.ConfirmPassword = "secrets"

	username := "rijik"
	data.Username = username

	_, status, err := modul.ChangePassword(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error updateting document: %v", err)
	} else {
		fmt.Println("Password berhasil diubah dengan username:", username)
	}
}

// test delete user
func TestDeleteUser(t *testing.T) {
	var data model.User
	data.Username = "rijik"

	status, err := modul.DeleteUser(mconn, "user", data)
	fmt.Println("Status", status)
	if err != nil {
		t.Errorf("Error deleting document: %v", err)
	} else {
		fmt.Println("Delete user" + data.Username + "success")
	}
}

func TestGetUserFromID(t *testing.T) {
	id := "656bf30c733cf24a0f73d0a8"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	}

	anu, err := modul.GetUserFromID(mconn, "user", ID)
	if err != nil {
		t.Errorf("Error getting user: %v", err)
		return
	}
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
	anu, _ := modul.GetUserFromEmail(mconn, "user", "rijik@gmail.com")
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
	var data model.Hp
	data.Title = "Vivo"
	data.Description = "vivo adalah sahabat"
	data.Deadline = "12/04/2023"
	// data.IsDone = false

	uid := "0040f398-1200-4f36-8332-6752ab3e55c0"

	id, err := modul.InsertHp(mconn, "hp", data, uid)
	if err != nil {
		t.Errorf("Error inserting Hp: %v", err)
	}
	fmt.Println(id)
}

func TestGetHpFromID(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("655c4408d06d3d2ddba5d1d7")
	anu, err := modul.GetHpFromID(mconn, "hp", id)
	if err != nil {
		t.Errorf("Error getting hp: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestGetHpFromUsername(t *testing.T) {
	anu, err := modul.GetHpFromUsername(mconn, "hp", "rijik")
	if err != nil {
		t.Errorf("Error getting hp: %v", err)
		return
	}
	fmt.Println(anu)
}

func TestUpdateHp(t *testing.T) {
	var data model.Hp
	data.Title = "Belajar Golang"
	data.Description = "Hari ini belajar golang"
	data.Deadline = "02/02/2021"

	id := "655c5047370b53741a9705d8"
	ID, err := primitive.ObjectIDFromHex(id)
	data.ID = ID
	if err != nil {
		fmt.Printf("Data tidak berhasil diubah")
	} else {

		_, status, err := modul.UpdateHp(mconn, "hp", data)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error updating hp with id: %v", err)
			return
		} else {
			fmt.Printf("Data berhasil diubah untuk id: %s\n", id)
		}
		fmt.Println(data)
	}
}

func TestDeleteHp(t *testing.T) {
	id := "655c4408d06d3d2ddba5d1d7"
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		t.Errorf("Error converting id to ObjectID: %v", err)
		return
	} else {

		status, err := modul.DeleteHp(mconn, "hp", ID)
		fmt.Println("Status", status)
		if err != nil {
			t.Errorf("Error deleting document: %v", err)
			return
		} else {
			fmt.Println("Delete success")
		}
	}
}
