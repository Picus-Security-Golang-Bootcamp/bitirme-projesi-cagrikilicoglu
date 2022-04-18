# Shopping Cart Api

![GitHub repo size](https://img.shields.io/github/repo-size/scottydocs/README-template.md)
![GitHub contributors](https://img.shields.io/github/contributors/scottydocs/README-template.md)
![GitHub stars](https://img.shields.io/github/stars/Picus-Security-Golang-Bootcamp/bitirme-projesi-cagrikilicoglu?style=social)
![GitHub forks](https://img.shields.io/github/forks/Picus-Security-Golang-Bootcamp/bitirme-projesi-cagrikilicoglu?style=social)
![Twitter Follow](https://img.shields.io/github/follow/cagrikilicoglu?style=social)

Shopping cart api simulates an e-commerce shopping cart that allows admin to create products and categories and user to do order products.

This api is designed as an example of a real e-commerce website so that all users can browse categories or products, registered users can complete shopping and admins can manipulate database.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed the latest version of Go.

## Installing Shopping Cart Api

To clone Shopping Cart Api to your repo follow these steps:

```
git clone https://github.com/Picus-Security-Golang-Bootcamp/bitirme-projesi-cagrikilicoglu.git
```

## Repository Overview

```bash
├── README.md
├── cmd
│   └── main.go
├── docs
│   └── shop.yml
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── cart.go
│   │   ├── category.go
│   │   ├── item.go
│   │   ├── login.go
│   │   ├── order.go
│   │   ├── product.go
│   │   ├── stock.go
│   │   └── user.go
│   ├── httpErrors
│   │   └── httpErrors.go
│   └── models
│       ├── cart
│       │   ├── handler.go
│       │   ├── repo.go
│       │   ├── repo_test.go
│       │   └── serializer.go
│       ├── category
│       │   ├── csvService.go
│       │   ├── handler.go
│       │   ├── repo.go
│       │   └── serializer.go
│       ├── item
│       │   ├── repo.go
│       │   ├── serializer.go
│       │   └── service.go
│       ├── models.go
│       ├── order
│       │   ├── handler.go
│       │   ├── repo.go
│       │   └── serializer.go
│       ├── product
│       │   ├── csvService.go
│       │   ├── handler.go
│       │   ├── repo.go
│       │   ├── repo_test.go
│       │   └── serializer.go
│       ├── response
│       │   └── response.go
│       └── user
│           ├── handler.go
│           ├── repo.go
│           ├── repo_test.go
│           └── serializer.go
├── pkg
│   ├── auth
│   │   └── auth.go
│   ├── config
│   │   ├── config.go
│   │   ├── local.yaml
│   │   └── production.yaml
│   ├── database
│   │   └── database.go
│   ├── graceful
│   │   └── shutdown.go
│   ├── jwtHelper
│   │   └── jwtHelper.go
│   ├── logging
│   │   ├── ginLogger.go
│   │   └── zapLogger.go
│   ├── middleware
│   │   └── middleware.go
│   └── pagination
│       └── pagination.go
└── test_file
    ├── categories.csv
    └── products.csv

```

## Features

- Authentication / Authorization
- JWT middleware for authentication
- File upload
- Database seed
- Pagination with Limit and Offset using GORM (Golang ORM framework)
- CRUD operations on products, categories, users, carts, orders
- Logging middlewares

## Getting Started

Before running the project you must set up your database credentials in the configuration file. Please fill DBConfig with your own dsn configuration including, DB_HOST, DB_PORT, DB_USERNAME, DB_NAME and DB_PASSWORD.

Also the default app environment set in the .env file is local, you can change it for production purposes, please configure your database accordingly.

## Using Shopping Cart Api

To use Shopping Cart Api, follow these steps:

The app provides the following endpoints:

#### Product

- `GET /api/v1/shopping-cart-api/products/` : list all the products with pagination parameters supplied by the user. If no pagination parameters are supplied, the endpoint uses defaults. ..

Example request: `GET /api/v1/shopping-cart-api/products/?page=3&pageSize=5`
requests the third page of all the products ordered by name and divided by groups of five.

```
`GET /api/v1/shopping-cart-api/products/sku/{sku}` : list the product with product SKU parameter.

Example request: `GET /api/v1/shopping-cart-api/products/sku/213DS`
requests product with SKU "213DS".

```

```
`GET /api/v1/shopping-cart-api/products/id/{id}` : list the product with product ID parameter.

Example request: `GET /api/v1/shopping-cart-api/products/sku/0f60fc10-4bed-4fcd-a5fe-d064a1a915cc`
requests product with ID "0f60fc10-4bed-4fcd-a5fe-d064a1a915cc".

```

```
`GET /api/v1/shopping-cart-api/products?name=?` : list the product/s with name parameter.

Example request: `GET /api/v1/shopping-cart-api/products?name=MacBook`
requests product/s whose name includes "MacBook". The search is elastic.

```

```
`POST /api/v1/shopping-cart-api/products/create` : creates a product supplied in the request body. The endpoint is only authorized for admin. Authorization token must be provided in the request header.

Example request: `POST /api/v1/shopping-cart-api/products/create`
requests body: {
    "categoryName": "Shoes",
    "name": "Nike Air Force 1",
    "price": 76,
    "stock": {
        "number": 20,
        "sku": "213DS",
    }
}
```

```
`POST /api/v1/shopping-cart-api/products/upload` : creates products from a csv file uploaded in the request body as a form file. The endpoint is only authorized for admin. Authorization token must be provided in the request header.
```

```
`PUT /api/v1/shopping-cart-api/products/update/sku/{sku}` : updates a product supplied in the request body. The endpoint is only authorized for admin. Authorization token must be provided in the request header.

Example request: `/api/v1/shopping-cart-api/products/update/sku/213DS`
requests body: {
    "categoryName": "Shoes",
    "name": "Nike Air Force 1",
    "price": 120,
    "stock": {
        "number": 50,
        "sku": "213DS",
    }
}
```

```
`DELETE /api/v1/shopping-cart-api/products/delete/sku/{sku}` : deletes a product with SKU parameter. The endpoint is only authorized for admin. Authorization token must be provided in the request header.

Example request: `DELETE /api/v1/shopping-cart-api/products/delete/sku/213DS`
```

#### Category

```
`GET /api/v1/shopping-cart-api/categories` : list all the categories with pagination parameters supplied by the user. If no pagination parameters are supplied, the endpoint uses defaults.

Example request: `GET /api/v1/shopping-cart-api/categories/?page=3&pageSize=5`
requests the third page of all the products ordered by name and divided by groups of five.
```

```
`GET /api/v1/shopping-cart-api/categories/{name}` : list the products of a category by catgory name input.

Example request: `GET /api/v1/shopping-cart-api/categories/Shoes`
requests product/s whose category is Shoes.
```

```
`POST /api/v1/shopping-cart-api/categories/create` : creates a category supplied in the request body. The endpoint is only authorized for admin. Authorization token must be provided in the request header.

Example request: `POST /api/v1/shopping-cart-api/categories/create`
requests body: {
    "Name":"Shoes",
    "Description": "ShoesDescription"
    }
}
```

```
`POST /api/v1/shopping-cart-api/categories/upload` : creates categories from a csv file uploaded in the request body as a form file. The endpoint is only authorized for admin. Authorization token must be provided in the request header.
```

#### User

```
`POST /api/v1/shopping-cart-api/signup` : registers a user supplied in the request body and logins the new user to the system by generating authorization tokens.

Example request: `POST /api/v1/shopping-cart-api/signup`
requests body: {
  "email": "test@mail.com",
  "password": "testPassword",
  "firstName": "Tester",
  "lastName": "Tested",
  "zipCode": "10001"
}
```

```
`POST /api/v1/shopping-cart-api/login` : logins a user whose credentials are supplied in the request body. Generate authorization tokens for the user

Example request: `POST /api/v1/shopping-cart-api/signup`
requests body: {
  "Email": "test@mail.com",
  "Password": "testPassword"
}
```

```
`POST /api/v1/shopping-cart-api/refresh` : the user gets two tokens when login is successful. Access token should be renewed with refresh token for certain intervals. The endpoint creates new tokens by existing refresh token. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `POST /api/v1/shopping-cart-api/refresh`
requests authorized user's tokens to refresh.
```

#### Cart

```
`GET /api/v1/shopping-cart-api/cart/` : shows the cart of the current user. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `GET /api/v1/shopping-cart-api/cart/`
requests authorized user's cart.
```

```
`POST /api/v1/shopping-cart-api/products/cart/add/sku/{sku}/quantity/{quantity}` : adds a product to the cart with SKU and quantity parameters. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `POST /api/v1/shopping-cart-api/products/cart/add/sku/12DSA/quantity/1`
requests adding the product with SKU 12DSA of quantity 1 to the authorized user's cart.
```

```
`PUT /api/v1/shopping-cart-api/products/cart/update/sku/{sku}/quantity/{quantity}` : updates the quantity of a product that is already in the cart with SKU and quantity parameters. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `PUT /api/v1/shopping-cart-api/products/cart/update/sku/12DSA/quantity/2`
requests updating the quantity of the product with SKU 12DSA to 2 in the authorized user's cart.
```

```
`DELETE /api/v1/shopping-cart-api/products/cart/delete/sku/{sku}` : deletes a product from the with SKU parameter. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `DELETE /api/v1/shopping-cart-api/products/cart/delete/sku/12DSA`
requests deleting the product with SKU 12DSA from the authorized user's cart.
```

#### Order

```
`POST /api/v1/shopping-cart-api/order` : orders products currently in the user's cart. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `POST /api/v1/shopping-cart-api/order`
requests ordering all the items in the authorized user's cart.
```

```
`DELETE /api/v1/shopping-cart-api/order/id/{id}/cancel` : cancels the order that is placed before with ID parameter. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `DELETE /api/v1/shopping-cart-api/order/id/82518cab-e9b0-4121-a51e-66e266b279s1/cancel`
request canceling the order with the ID 82518cab-e9b0-4121-a51e-66e266b279s1 of authorized user.
```

```
`GET /api/v1/shopping-cart-api/order/history` : gets all the order history. The endpoint is only authorized for admin and user. Authorization token must be provided in the request header.

Example request: `GET /api/v1/shopping-cart-api/order/history`
requests order history of authorized user.
```

## Tool set

- Go
- Gorm
- Gin
- Viper
- jwt-go
- Postgresql
- Swagger
- zap
- sqlMock

## Contributors

Thanks to the following people who have contributed to this project:

- [@cagrikilicoglu](https://github.com/cagrikilicoglu) 📖

## Links

[project repository](https://github.com/Picus-Security-Golang-Bootcamp/bitirme-projesi-cagrikilicoglu)

## License

This project uses the following license: [MIT](https://opensource.org/licenses/MIT).
