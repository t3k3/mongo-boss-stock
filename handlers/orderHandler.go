package handlers

import (
	"bytes"
	"context"
	"log"
	"math"
	"net/http"
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

	order.CreatedAt = time.Now()

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

	//Burada order içindeki OrderProducts nesnelerinin ID.Hex'leri stok düşümü için
	//Yazdığımız fonksiyona adet sayısının hesaplanmasıyla birlikte gönderiliyor.
	//Burada local bir http request atılıyor. Best Practice değil ancak şimdilik çalışıyor.
	//salesQty yapılan satış adedini saklarken stockQty ürünün güncel stoğunu tutmaktadır.

	for _, v := range order.OrderProducts {

		UpdateProductQty(v.ID.Hex(), v.StockQuantity-v.CartQuantity)
		//UpdateProductQty(v.ID.Hex(), v.SalesQty+v.SalesQty)

	}

	//TODO: Satış yapılan ürünün stoktan düşmesi için

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Order inserted successfully",
	})

}
func UpdateProductQty(stockQty string, salesQty int) {

	url := "http://localhost:3001/api/products/" + stockQty

	data := `{"stock_quantity":` + strconv.Itoa(int(salesQty)) + `}`

	var jsonStr = []byte(data)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

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
