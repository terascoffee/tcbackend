package rkbackend

import (
	"context"
	"errors"
	"fmt"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func UpdateOneDoc(db *mongo.Database, col string, filter, update interface{}) (err error) {
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		fmt.Printf("UpdateOneDoc: %v\n", err)
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified filter")
		return err
	}
	return
}

func DeleteOneDoc(db *mongo.Database, col string, filter bson.M) (err error) {
	cols := db.Collection(col)
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Printf("DeleteOneDoc: %v\n", err)
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("no data has been deleted with the specified filter")
		return err
	}
	return
}

// Admin
func InsertProduk(db *mongo.Database, col string, produkdata Produk) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, produkdata)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
	}
	return insertedID, err
}

func InsertTransaksi(db *mongo.Database, col string, transaksidata Transaksi) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, transaksidata)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
	}
	return insertedID, err
}

func GetAllDataProduk(db *mongo.Database, col string) (produklist []Produk) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &produklist)
	if err != nil {
		fmt.Println(err)
	}
	return produklist
}

func GetAllDataTransaksi(db *mongo.Database, col string) (transaksilist []Produk) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &transaksilist)
	if err != nil {
		fmt.Println(err)
	}
	return transaksilist
}

func InsertUser(db *mongo.Database, collection string, userdata Admin) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Username : " + userdata.Username + "\nPassword : " + userdata.Password
}

func UpdateProduk(db *mongo.Database, col string, produk Produk) (produks Produk, status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": produk.ID}
	update := bson.M{
		"$set": bson.M{
			"nama":      produk.Nama,
			"harga":     produk.Harga,
			"deskripsi": produk.Deskripsi,
			"stok":      produk.Stok,
		},
	}

	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return produks, false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil diupdate")
		return produks, false, err
	}

	err = cols.FindOne(context.Background(), filter).Decode(&produks)
	if err != nil {
		return produks, false, err
	}

	return produks, true, nil
}

func DeleteProduk(db *mongo.Database, col, nama string) (status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"nama": nama}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("Data tidak berhasil dihapus")
		return false, err
	}
	return true, nil
}

func GetProdukFromID(db *mongo.Database, col string, _id primitive.ObjectID) (produklist Produk, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	err = cols.FindOne(context.Background(), filter).Decode(&produklist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			fmt.Println("no data found for ID", _id)
		} else {
			fmt.Println("error retrieving data for ID", _id, ":", err.Error())
		}
	}

	return produklist, nil
}

func GetProdukFromName(db *mongo.Database, col string, nama string) (produklist []Produk, err error) {
	cols := db.Collection(col)
	filter := bson.M{"nama": nama}

	cursor, err := cols.Find(context.Background(), filter)
	if err != nil {
		fmt.Println("Error GetProdukFromName in colection", col, ":", err)
		return nil, err
	}

	err = cursor.All(context.Background(), &produklist)
	if err != nil {
		fmt.Println(err)
	}

	return produklist, nil
}
