# buyer-seller-app


This repo consists of a seller service is a web application that provides API for managing sellers and searching for products 

## Prerequisites
- Docker 
- Docker Compose

This service uses **Mysql** as its datastore.

## Getting Started

```shell
git clone https://github.com/ganesh-sai/buyer-seller-app.git
```

Start the application using docker compose; it brings up all the necessary dependencies and runs app inside a docker container
```shell
docker-compose -f docker-compose.yml up
```

There exists a [init.sql](./init.sql) which contains the schema. Execute it on the db using your preferred choice, like dbBeaver, mysql workbench.

Navigate to `seller-service`
```shell
cd buyer-seller-app/seller-service
```
In the above directory, you will find code for the app that it is running.

The app will run on port 8080
```
http://localhost:8080
```


The below are the endpoints that are provided by the app

## Create a Product [API](./seller-service/handlers/product_handler.go)

- Endpoint: `POST /api/v1/product`
- Input:
    ``` 
  {
  "sellerId": 1,
  "productName": "Product Name",
  "price": 10.0,
  "quantity": 5
  }
  ```
- Output: returns a product id that is saved to the db
  ```
  {
      "id": 1
  }
  ```

## Create a Seller [API](./seller-service/handlers/seller_handler.go)
- Endpoint: `POST /api/v1/seller`
- Input: 
    ```
  {
    "name": "Seller Name",
    "location": "Seller Location"
  }
  ```
- Output: returns a plain string back to the user with seller id in it.
    ```
    Seller created with ID: 1
  ```
  
## Search Products [API](./seller-service/handlers/products_search_handler.go)
- Endpoint: `GET /api/v1/product/search`
- Query Parameters: 
  - `productName` (optional): Product name for filtering products
  - `desiredQty` (optional): Desired quantity for filtering products
  - `location` (optional): Location for filtering products
  - `minPrice` (optional): Minimum price for filtering products
  - `maxPrice` (optional): Maximum price for filtering products
  - `sortBy` (optional): Field to sort the products (available:  "price", "productName", "sellerId", "productId")
  - `page` (optional): Page number for pagination
  - `perPage` (optional): Number of products per page 


- Output: 
 ```
[
  {
    "id": 1,
    "sellerId": 1,
    "productName": "Product Name",
    "price": 10.0,
    "quantity": 5
  },
  {
    "id": 2,
    "sellerId": 1,
    "productName": "Another Product",
    "price": 20.0,
    "quantity": 3
  },
  ...
]
```