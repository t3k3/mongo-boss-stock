package handlers

import (
	"context"
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nfnt/resize"
	"github.com/t3k3/mongo-boss-stock/config"
	"github.com/t3k3/mongo-boss-stock/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllProducts(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var products []models.Product

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
					"product": bson.M{
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

	total, _ := productCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := productCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Products Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var product models.Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      products,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func UploadFile(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	// Check for errors:

	filePath := ""

	trimmedFileName := strings.Join(strings.Fields(file.Filename), "")

	if err == nil {

		//Save file inside uploads folder under current working directory:
		c.SaveFile(file, fmt.Sprintf("./uploads/%s", trimmedFileName))

		filePath = c.BaseURL() + "/static/thumbs/" + trimmedFileName

		//RESİZE START
		// open "test.jpg"
		file, err := os.Open("./uploads/" + trimmedFileName)
		if err != nil {

			log.Fatal(err)
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {

			log.Fatal(err)
		}
		file.Close()
		// resize to width 1000 using Lanczos resampling
		// and preserve aspect ratio
		m := resize.Thumbnail(80, 100, img, resize.Lanczos3)

		out, err := os.Create("./uploads/thumbs/" + trimmedFileName)
		if err != nil {

			log.Fatal(err)
		}
		defer out.Close()

		// write new image to file
		jpeg.Encode(out, m, nil)
		//RESİZE END

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    filePath,
		"success": true,
		"error":   err,
	})
}

func GetProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var product models.Product
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := productCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&product)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    product,
		"success": true,
	})
}

func AddProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	product := new(models.Product)

	product.CreatedAt = time.Now()
	product.DeletedAt = time.Now()

	if err := c.BodyParser(product); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := productCollection.InsertOne(ctx, product)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Product inserted successfully",
	})

}

func UpdateProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil {
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
			"message": "Product not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": product,
	}
	_, err = productCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product updated successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err,
		})
	}
	_, err = productCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}
