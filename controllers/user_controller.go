package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"helloworld/configs"
	"helloworld/models"
	"helloworld/redis"
	"helloworld/responses"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

/* Creat User Interface
accept Parameter:
	UserId    int
	FirstName string
	LastName  string
	Age		  int
	Gender	  string
*/
func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "validate body error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "validate required field error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.User{
		UserId:    user.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		Gender:    user.Gender,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result, "message": "User create successfully!"}})
}

// Get All User Interface
func GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	// connect redis server
	redisInstant := redis.ConnectRedisServer()

	cacheRes, err := redisInstant.Get(ctx, "GetAllUser").Result()
	if len(cacheRes) != 0 {
		return c.Status(http.StatusOK).JSON(
			responses.UserResponse{Status: http.StatusOK, Message: "Success, return by cache", Data: &fiber.Map{"data": cacheRes}})
	}

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		var err = results.Decode(&singleUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	// Set redis cache for GetAllUser response
	userJson, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Convert Get All User map to JSON")
	}

	cacheErr := redisInstant.Set(ctx, "GetAllUser", userJson, 10*time.Second).Err()
	if cacheErr != nil {
		fmt.Println(cacheErr)
	} else {
		fmt.Println("GetAllUser data set to redis cache!")
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

// Get Singel User Interface
func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	userIdInt, convErr := strconv.Atoi(userId)
	if convErr != nil {
	}

	err := userCollection.FindOne(ctx, bson.M{"userid": userIdInt}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

// Update User Interface
func UpdateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user models.User
	defer cancel()

	userIdInt, convErr := strconv.Atoi(userId)
	if convErr != nil {
	}

	// validate request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	updateUser := bson.M{"firstname": user.FirstName, "lastname": user.LastName, "age": user.Age, "gender": user.Gender} // Should stay in one line

	updateResult, err := userCollection.UpdateOne(ctx, bson.M{"userid": userIdInt}, bson.M{"$set": updateUser})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "Update error", Data: &fiber.Map{"data": err.Error()}})
	}

	if updateResult.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"userid": userIdInt}).Decode(&updateUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "Multiple Match error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updateUser, "message": "User update successfully!"}})
}

// Delete User Interface
func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	userIdInt, convErr := strconv.Atoi(userId)
	if convErr != nil {
	}

	deleteResult, err := userCollection.DeleteOne(ctx, bson.M{"userid": userIdInt})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "UserId Not FOUND!"}})
	}

	if deleteResult.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"message": "User delete successfully"}})
}

// Get User Count
func CountUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"Total User Count": count}})
}
