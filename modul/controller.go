package modul

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aiteung/atdb"
	"github.com/badoux/checkmail"
	model "github.com/wegotour/be_p3/model"
)

func MongoConnect(MONGOCONNSTRINGENV, dbname string) *mongo.Database {
	var DBmongoinfo = atdb.DBInfo{
		// DBString: "mongodb+srv://admin:admin@projectexp.pa7k8.gcp.mongodb.net", //os.Getenv(MONGOCONNSTRINGENV),
		DBString: os.Getenv(MONGOCONNSTRINGENV),
		DBName:   dbname,
	}
	return atdb.MongoConnect(DBmongoinfo)
}

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return insertedID, err
}

// user
func Register(db *mongo.Database, col string, userdata model.User) error {
	if userdata.Username == "" || userdata.Password == "" || userdata.Email == "" {
		return fmt.Errorf("data tidak lengkap")
	}

	// Periksa apakah email valid
	if err := checkmail.ValidateFormat(userdata.Email); err != nil {
		return fmt.Errorf("email tidak valid")
	}

	// Periksa apakah email dan username sudah terdaftar
	userExists, _ := GetUserFromEmail(db, col, userdata.Email)
	if userExists.Email != "" {
		return fmt.Errorf("email sudah terdaftar")
	}
	userExists, _ = GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username != "" {
		return fmt.Errorf("username sudah terdaftar")
	}

	// Periksa apakah password memenuhi syarat
	if len(userdata.Password) < 6 {
		return fmt.Errorf("password minimal 6 karakter")
	}
	if strings.Contains(userdata.Password, " ") {
		return fmt.Errorf("password tidak boleh mengandung spasi")
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		return fmt.Errorf("username tidak boleh mengandung spasi")
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	user := bson.M{
		"_id":      primitive.NewObjectID(),
		"email":    userdata.Email,
		"username": userdata.Username,
		"password": hash,
		"role":     "user",
	}
	_, err := InsertOneDoc(db, col, user)
	if err != nil {
		return fmt.Errorf("SignUp: %v", err)
	}
	return nil
}

func LogIn(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Password == "" {
		err = fmt.Errorf("data tidak lengkap")
		return user, false, err
	}

	// Periksa apakah pengguna dengan username tertentu ada
	userExists, _ := GetUserFromUsername(db, col, userdata.Username)
	if userExists.Username == "" {
		err = fmt.Errorf("username tidak ditemukan")
		return user, false, err
	}

	// Periksa apakah kata sandi benar
	if !CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("password salah")
		return user, false, err
	}
	return userExists, true, nil
}

func UpdateUser(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	if userdata.Username == "" || userdata.Email == "" {
		err = fmt.Errorf("data tidak boleh kosong")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	existingUser, err := GetUserFromID(db, col, userdata.ID)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah data yang akan diupdate sama dengan data yang sudah ada
	if userdata.Username == existingUser.Username && userdata.Email == existingUser.Email {
		err = fmt.Errorf("data yang ingin diupdate tidak boleh sama")
		return user, false, err
	}

	checkmail.ValidateFormat(userdata.Email)
	if err != nil {
		err = fmt.Errorf("email tidak valid")
		return user, false, err
	}

	// Periksa apakah username memenuhi syarat
	if strings.Contains(userdata.Username, " ") {
		err = fmt.Errorf("username tidak boleh mengandung spasi")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	filter := bson.M{"_id": userdata.ID}
	update := bson.M{
		"$set": bson.M{
			"email":    userdata.Email,
			"username": userdata.Username,
			"password": hash,
			"role":     "user",
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("data tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func ChangePassword(db *mongo.Database, col string, userdata model.User) (user model.User, status bool, err error) {
	// Periksa apakah pengguna dengan username tertentu ada
	userExists, err := GetUserFromUsername(db, col, userdata.Username)
	if err != nil {
		return user, false, err
	}

	// Periksa apakah password memenuhi syarat
	if userdata.Password == "" {
		err = fmt.Errorf("password tidak boleh kosong")
		return user, false, err
	}
	if len(userdata.Password) < 6 {
		err = fmt.Errorf("password minimal 6 karakter")
		return user, false, err
	}
	if strings.Contains(userdata.Password, " ") {
		err = fmt.Errorf("password tidak boleh mengandung spasi")
		return user, false, err
	}

	// Periksa apakah password sama dengan password lama
	if CheckPasswordHash(userdata.Password, userExists.Password) {
		err = fmt.Errorf("password tidak boleh sama")
		return user, false, err
	}

	// Simpan pengguna ke basis data
	hash, _ := HashPassword(userdata.Password)
	userExists.Password = hash
	filter := bson.M{"username": userdata.Username}
	update := bson.M{
		"$set": bson.M{
			"email":    userdata.Email,
			"username": userdata.Username,
			"password": userExists.Password,
			"role":     "user",
		},
	}
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return user, false, err
	}
	if result.ModifiedCount == 0 {
		err = fmt.Errorf("PAssword tidak berhasil diupdate")
		return user, false, err
	}
	return user, true, nil
}

func DeleteUser(db *mongo.Database, col string, username string) error {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("data tidak berhasil dihapus")
	}

	return nil
}

func GetUserFromID(db *mongo.Database, col string, _id primitive.ObjectID) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromID: %v\n", err)
	}
	return user, nil
}

func GetUserFromUsername(db *mongo.Database, col string, username string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"username": username}
	err = cols.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromUsername: %v\n", err)
		return user, err
	}
	return user, nil
}

func GetUserFromEmail(db *mongo.Database, col string, email string) (user model.User, err error) {
	cols := db.Collection(col)
	filter := bson.M{"email": email}
	err = cols.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Printf("GetUserFromEmail: %v\n", err)
	}
	return user, nil
}

func GetAllUser(db *mongo.Database, col string) (userlist []model.User, err error) {
	ctx := context.TODO()
	cols := db.Collection(col)
	filter := bson.M{}

	cur, err := cols.Find(ctx, filter)
	if err != nil {
		fmt.Println("Error GetAllUser in colection", col, ":", err)
		return nil, err
	}

	// defer cur.Close(ctx)
	defer func() {
		if cerr := cur.Close(ctx); cerr != nil {
			fmt.Println("Error closing cursor:", cerr)
		}
	}()

	err = cur.All(context.TODO(), &userlist)
	if err != nil {
		fmt.Println("Error reading documents:", err)
		return nil, err
	}

	return userlist, nil
}

// hp
func InsertHp(db *mongo.Database, col string, hp model.Hp) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, hp)
	if err != nil {
		fmt.Printf("InsertHp: %v\n", err)
	}
	return insertedID, err
}

func GetHpFromID(db *mongo.Database, col string, id primitive.ObjectID) (hp model.Hp) {
	cols := db.Collection(col)
	filter := bson.M{"_id": id}
	err := cols.FindOne(context.Background(), filter).Decode(&hp)
	if err != nil {
		fmt.Printf("GetHpFromID: %v\n", err)
	}
	return hp
}

func GetHpList(db *mongo.Database, col string) (hp []model.Hp, err error) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetHpList in colection", col, ":", err)
		return nil, err
	}
	err = cursor.All(context.Background(), &hp)
	if err != nil {
		fmt.Println(err)
	}
	return hp, nil
}

func UpdateHp(db *mongo.Database, col string, hp model.Hp) (hps model.Hp, status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": hp.ID}
	update := bson.M{
		"$set": bson.M{
			"title":               hp.Title,
			"description":         hp.Description,
			"deadline":            hp.Deadline,
			"timestamp.updatedat": time.Now(),
		},
		"$setOnInsert": bson.M{
			"timestamp.createdat": hp.TimeStamp.CreatedAt,
		},
	}

	options := options.Update().SetUpsert(true)

	result, err := cols.UpdateOne(context.Background(), filter, update, options)
	if err != nil {
		return hps, false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		err = fmt.Errorf("data tidak berhasil diupdate")
		return hps, false, err
	}

	err = cols.FindOne(context.Background(), filter).Decode(&hps)
	if err != nil {
		return hps, false, err
	}

	return hps, true, nil
}

func DeleteHp(db *mongo.Database, col string, _id primitive.ObjectID) (status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("data tidak berhasil dihapus")
		return false, err
	}
	return true, nil
}
