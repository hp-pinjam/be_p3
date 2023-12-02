package be_p3

import (
	"fmt"
	"testing"

	model "github.com/wegotour/be_p3/model"
	modul "github.com/wegotour/be_p3/modul"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mconn = SetConnection("MONGOSTRING", "wegotour")

// user
func TestRegister(t *testing.T) {
	var data model.User
	data.ID = primitive.NewObjectID()
	data.Email = "dapskuy@gmail.com"
	data.Username = "dapskuy"
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
	userdata.Username = "dapskuy"
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
	data.Email = "dapskuy@gmail.com"
	data.Username = "dapskuy"
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
	data.Email = "dapskuy@gmail.com" // email tidak diubah
	data.Username = "dapskuy"        // username tidak diubah
	data.Role = "admin"              // role tidak diubah

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
	username := "dapskuy"

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
	anu, err := modul.GetUserFromUsername(mconn, "user", "dapskuy")
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

// ticket
func TestInsertTicket(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "wegotour")
	var ticketdata model.Ticket
	ticketdata.Title = "Perjalanan"
	ticketdata.Description = "pergi ke bali"
	ticketdata.IsDone = true

	nama, err := modul.InsertTicket(mconn, "ticket", ticketdata)
	if err != nil {
		t.Errorf("Error inserting ticket: %v", err)
	}
	fmt.Println(nama)
}

func TestGetTicketFromID(t *testing.T) {
	mconn := SetConnection("MONGOSTRING", "wegotour")
	id, _ := primitive.ObjectIDFromHex("6548bce6c31c8ec3f02fa11d")
	anu := modul.GetTicketFromID(mconn, "ticket", id)
	fmt.Println(anu)
}

func TestGetTicketList(t *testing.T) {
	anu, err := modul.GetTicketList(mconn, "ticket")
	if err != nil {
		t.Errorf("Error getting ticket: %v", err)
		return
	}
	fmt.Println(anu)
}
