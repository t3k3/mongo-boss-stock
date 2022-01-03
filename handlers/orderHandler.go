package handlers

import (
	"context"
	"fmt"
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

func GetAllOrders(c *fiber.Ctx) error {
	orderCollection := config.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var orders []models.Order

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
					"order": bson.M{
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

	total, _ := orderCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := orderCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Orders Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var order models.Order
		cursor.Decode(&order)
		orders = append(orders, order)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      orders,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetOrder(c *fiber.Ctx) error {
	orderCollection := config.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var order models.Order
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := orderCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&order)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    order,
		"success": true,
	})
}

func AddOrder(c *fiber.Ctx) error {
	orderCollection := config.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(models.Order)

	if err := c.BodyParser(order); err != nil {
		log.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	result, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to insert",
			"error":   err,
		})
	}
	//TODO: Satış yapılan ürünün stoktan düşmesi için

	//m := make(map[string]int)
	var s1 [10]string
	for v, p := range order.OrderProducts {
		s1[v] = p.ID.Hex()
	}
	fmt.Println(s1)

	//TODO: Satış yapılan ürünün stoktan düşmesi için

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Order inserted successfully",
	})

}

func UpdateOrder(c *fiber.Ctx) error {
	orderCollection := config.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	order := new(models.Order)

	if err := c.BodyParser(order); err != nil {
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
			"message": "Order not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": order,
	}
	_, err = orderCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Order updated successfully",
	})
}

func DeleteOrder(c *fiber.Ctx) error {
	orderCollection := config.MI.DB.Collection("orders")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Order not found",
			"error":   err,
		})
	}
	_, err = orderCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Order failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Order deleted successfully",
	})
}
