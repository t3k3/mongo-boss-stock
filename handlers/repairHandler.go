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

func GetAllRepairs(c *fiber.Ctx) error {
	repairCollection := config.MI.DB.Collection("repairs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var repairs []models.Repair

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
					"repair": bson.M{
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

	total, _ := repairCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := repairCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Repairs Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var repair models.Repair
		cursor.Decode(&repair)
		repairs = append(repairs, repair)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      repairs,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetRepair(c *fiber.Ctx) error {
	repairCollection := config.MI.DB.Collection("repairs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var repair models.Repair
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := repairCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Repair Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&repair)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Repair Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    repair,
		"success": true,
	})
}

func AddRepair(c *fiber.Ctx) error {
	repairCollection := config.MI.DB.Collection("repairs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repair := new(models.Repair)

	if err := c.BodyParser(repair); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := repairCollection.InsertOne(ctx, repair)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Repair failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Repair inserted successfully",
	})

}

func UpdateRepair(c *fiber.Ctx) error {
	repairCollection := config.MI.DB.Collection("repairs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	repair := new(models.Repair)

	if err := c.BodyParser(repair); err != nil {
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
			"message": "Repair not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": repair,
	}
	_, err = repairCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Repair failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Repair updated successfully",
	})
}

func DeleteRepair(c *fiber.Ctx) error {
	repairCollection := config.MI.DB.Collection("repairs")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Repair not found",
			"error":   err,
		})
	}
	_, err = repairCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Repair failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Repair deleted successfully",
	})
}
