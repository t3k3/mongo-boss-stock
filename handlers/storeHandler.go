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

func GetAllStores(c *fiber.Ctx) error {
	storeCollection := config.MI.DB.Collection("stores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var stores []models.Store

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
					"store": bson.M{
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

	total, _ := storeCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := storeCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Stores Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var store models.Store
		cursor.Decode(&store)
		stores = append(stores, store)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      stores,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetStore(c *fiber.Ctx) error {
	storeCollection := config.MI.DB.Collection("stores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var store models.Store
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := storeCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Store Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&store)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Store Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    store,
		"success": true,
	})
}

func AddStore(c *fiber.Ctx) error {
	storeCollection := config.MI.DB.Collection("stores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	store := new(models.Store)

	if err := c.BodyParser(store); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := storeCollection.InsertOne(ctx, store)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Store failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Store inserted successfully",
	})

}

func UpdateStore(c *fiber.Ctx) error {
	storeCollection := config.MI.DB.Collection("stores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	store := new(models.Store)

	if err := c.BodyParser(store); err != nil {
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
			"message": "Store not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": store,
	}
	_, err = storeCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Store failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Store updated successfully",
	})
}

func DeleteStore(c *fiber.Ctx) error {
	storeCollection := config.MI.DB.Collection("stores")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Store not found",
			"error":   err,
		})
	}
	_, err = storeCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Store failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Store deleted successfully",
	})
}
