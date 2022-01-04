
# Boss-Stock Fiber (Go) and MongoDB REST API

An API built with Fiber and MongoDB.

## Installation

```bash
go mod init
go get -u github.com/gofiber/fiber/v2
go get go.mongodb.org/mongo-driver/mongo
go get github.com/joho/godotenv
```

## API Usage

#### Get all products

```http
  GET /products
```

| Optional Parameter | Type     | Description                                        | Example       |
| :----------------- | :------- | :------------------------------------------------- | :------------ |
| `s`                | `string` | Serach Product Name or Products                   | ?s=exmpletext |
| `page`             | `int`	| Page number. Default: 1                            | ?page=2       |
| `limit`	     | `int`    | Limit number of products per page. Default: 10 | ?limit=20     |

#### Get products

```http
  GET /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of products to fetch |


#### Add products

```http
  POST /products
```


```json
{
    "name": "Apple",
    "detail": "Apple AirPods 2",
    "price": 25.80,
    "quantity":20,
    "barcode": 1003,
    "store_id": 2,
    "category_name": "Telefon",
    "entry_price": 20.80,
    "kdv": 18.0
}

```

#### Update catchphrase

```http
  PATCH /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of product to update |


```json
{
    "name": "Apple",
    "detail": "Apple AirPods 3", //Update AirPods Model 2 to 3
    "price": 25.80,
    "quantity":20,
    "barcode": 1003,
    "store_id": 2,
    "category_name": "Telefon",
    "entry_price": 20.80,
    "kdv": 18.0
}

// All nine fields are optional when updating

```

#### Remove product

```http
  DELETE /products/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of products to delete |


## Thanks

[Tutorial](https://dev.to/mikefmeyer/build-a-go-rest-api-with-fiber-and-mongodb-44og) https://dev.to/mikefmeyer/build-a-go-rest-api-with-fiber-and-mongodb-44og

