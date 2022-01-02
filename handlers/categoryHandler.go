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

func GetAllCategories(c *fiber.Ctx) error {
	categoryCollection := config.MI.DB.Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var categories []models.Category

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
					"category": bson.M{
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

	total, _ := categoryCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := categoryCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Categories Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var category models.Category
		cursor.Decode(&category)
		categories = append(categories, category)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      categories,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetCategory(c *fiber.Ctx) error {
	categoryCollection := config.MI.DB.Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var category models.Category
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := categoryCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Category Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&category)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Category Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    category,
		"success": true,
	})
}

func AddCategory(c *fiber.Ctx) error {
	categoryCollection := config.MI.DB.Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := categoryCollection.InsertOne(ctx, category)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Category failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Category inserted successfully",
	})

}

func UpdateCategory(c *fiber.Ctx) error {
	categoryCollection := config.MI.DB.Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	category := new(models.Category)

	if err := c.BodyParser(category); err != nil {
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
			"message": "Category not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": category,
	}
	_, err = categoryCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Category failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Category updated successfully",
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	categoryCollection := config.MI.DB.Collection("categories")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Category not found",
			"error":   err,
		})
	}
	_, err = categoryCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Category failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Category deleted successfully",
	})
}
