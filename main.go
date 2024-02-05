package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var kelasColl *mongo.Collection
var schoolColl *mongo.Collection

func getSchoolList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := schoolColl.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var schoolList []School
	err = cursor.All(ctx, &schoolList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schoolList)
}

func getKelasList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := kelasColl.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var kelasList []Kelas
	err = cursor.All(ctx, &kelasList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kelasList)
}

func createSchool(c *gin.Context) {
	var newSchool School

	if err := c.ShouldBindJSON(&newSchool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSchool = *NewSchool(newSchool.Name)

	result, err := schoolColl.InsertOne(context.Background(), newSchool)
	if err != nil {
		log.Println("Error inserting school:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("School inserted successfully:", result.InsertedID)
	c.JSON(http.StatusCreated, newSchool)
}

func createKelas(c *gin.Context) {
	var newKelas Kelas

	if err := c.ShouldBindJSON(&newKelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newKelas = *NewKelas(newKelas.Sekolah, newKelas.Name)

	result, err := kelasColl.InsertOne(context.Background(), newKelas)
	if err != nil {
		log.Println("Error inserting kelas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("kelas inserted successfully:", result.InsertedID)
	c.JSON(http.StatusCreated, newKelas)
}
func updateKelas(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedKelas Kelas
	if err := c.ShouldBindJSON(&updatedKelas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = kelasColl.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updatedKelas},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kelas updated successfully"})
}

func getKelasByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var kelas Kelas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = kelasColl.FindOne(ctx, bson.M{"_id": id}).Decode(&kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sekolah School
	err = schoolColl.FindOne(ctx, bson.M{"_id": kelas.Sekolah}).Decode(&sekolah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kelas)
}

func getKelasDetailByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var kelas Kelas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = kelasColl.FindOne(ctx, bson.M{"_id": id}).Decode(&kelas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var sekolah School
	err = schoolColl.FindOne(ctx, bson.M{"_id": kelas.Sekolah}).Decode(&sekolah)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kelas)
}

func getSchoolByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var school School
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = schoolColl.FindOne(ctx, bson.M{"_id": id}).Decode(&school)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, school)
}

func updateSchool(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedSchool School
	if err := c.ShouldBindJSON(&updatedSchool); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = schoolColl.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updatedSchool},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "School updated successfully"})
}

func main() {
	InitMongoDB()
	defer client.Disconnect(context.Background())
	router := gin.Default()

	router.GET("/school", getSchoolList)
	router.GET("/school/:id", getSchoolByID)
	router.POST("/school", createSchool)
	router.PUT("/school/:id", updateSchool)

	router.GET("/kelas", getKelasList)
	router.GET("/kelas/:id", getKelasByID)

	router.POST("/kelas", createKelas)
	router.PUT("/kelas/:id", updateKelas)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
