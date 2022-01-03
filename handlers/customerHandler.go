package handlers

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/config"
	"github.com/t3k3/mongo-boss-stock/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllCustomers(c *fiber.Ctx) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var customers []models.Customer

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"name": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"customer": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "10"))
	var limit int64 = int64(limitVal)

	total, _ := customerCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := customerCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "customers Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var customer models.Customer
		cursor.Decode(&customer)
		customers = append(customers, customer)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      customers,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetCustomer(c *fiber.Ctx) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var customer models.Customer
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := customerCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "customer Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Customer Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    customer,
		"success": true,
	})
}

func AddCustomer(c *fiber.Ctx) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	customer := new(models.Customer)

	if err := c.BodyParser(customer); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := customerCollection.InsertOne(ctx, customer)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Customer failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Customer inserted successfully",
	})

}

func UpdateCustomer(c *fiber.Ctx) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	customer := new(models.Customer)

	if err := c.BodyParser(customer); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Customer not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": customer,
	}
	_, err = customerCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Customer failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Customer updated successfully",
	})
}

func DeleteCustomer(c *fiber.Ctx) error {
	customerCollection := config.MI.DB.Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Customer not found",
			"error":   err,
		})
	}
	_, err = customerCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Customer failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Customer deleted successfully",
	})
}
