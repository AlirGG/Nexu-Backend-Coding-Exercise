Building, Running, and Testing the RESTful API
Introduction
This RESTful API is designed to provide information about car brands and their models from a MongoDB database. This document provides instructions for building, running, and testing the API.


Prerequisites
Before you can build and run the API, you must have the following software installed on your system:

Golang 1.15 or later
MongoDB 4.0 or later
Building the API


To build the API, please follow these steps:

Open a terminal window and navigate to the root directory of the API.
Run the following command to download the necessary dependencies:

go mod download


Run the following command to build the API binary:

go build -o api


Running the API
To run the API, please follow these steps:

Open a terminal window and navigate to the root directory of the API.
Run the following command to start the API:

go run main.go
The API will start running on port 8000. You can access the API using a web browser or a tool like Postman or Thunder Client in VSCode(used for internal tests).

Example of tests output in Thunder Client:

GET http://localhost:8000/models?greater=380000&lower=400000

response:

[
  {
    "id": 171,
    "name": "Wrangler",
    "average_price": 396757,
    "brand_name": "Jeep"
  },
  {
    "id": 194,
    "name": "CX9",
    "average_price": 383370,
    "brand_name": "Mazda"
  },
  {
    "id": 382,
    "name": "X3",
    "average_price": 398124,
    "brand_name": "BMW"
  },
  {
    "id": 737,
    "name": "SEI2",
    "average_price": 383714,
    "brand_name": "JAC"
  },
  {
    "id": 1526,
    "name": "Renegade",
    "average_price": 396920,
    "brand_name": "Jeep"
  },
  {
    "id": 1542,
    "name": "ELF 200",
    "average_price": 380933,
    "brand_name": "Isuzu"
  },
  {
    "id": 1547,
    "name": "Ram Promaster",
    "average_price": 389350,
    "brand_name": "RAM"
  }
]


GET http://localhost:8000/brands

response:

[
  {
    "average_price": 3342575,
    "id": 1,
    "nombre": "Bentley"
  },
  {
    "average_price": 209228.76,
    "id": 2,
    "nombre": "Dodge"
  },
  {
    "average_price": 0,
    "id": 3,
    "nombre": "Mastretta"
  },
  {
    "average_price": 0,
    "id": 4,
    "nombre": "Saab"
  },
  {
    "average_price": 630759.4666666667,
    "id": 5,
    "nombre": "Audi"
  },
  {
    "average_price": 224084.58620689655,
    "id": 6,
    "nombre": "Volkswagen"
  },
  {
    "average_price": 107028.17647058824,
    "id": 7,
    "nombre": "Fiat"
  },
  {
    "average_price": 176244.84210526315,
    "id": 8,
    "nombre": "Ford"
  },
  {
    "average_price": 854875.0909090909,
    "id": 9,
    "nombre": "Jaguar"
  },
  {
    "average_price": 205150,
    "id": 10,
    "nombre": "Mitsubishi"
  },
  {
    "average_price": 415593.25,
    "id": 11,
    "nombre": "Volvo"
  },
  {
    "average_price": 216107.35714285713,
    "id": 12,
    "nombre": "Hyundai"
  },
  {
    "average_price": 479268.22222222225,
    "id": 13,
    "nombre": "Cadillac"
  },
  {
    "average_price": 0,
    "id": 14,
    "nombre": "Lamborghini"
  },
  {
    "average_price": 823132,
    "id": 15,
    "nombre": "Land Rover"
  },
  {
    "average_price": 263356.90476190473,
    "id": 16,
    "nombre": "Toyota"
  },
  ...]


Automated Tests were contemplated but not finished in time, a test_controller.txt doc is made where the code for test_controller.go is implemented but it has some bugs to handle that cannot been handled in time.
Consider that code as automated test and use tools like Postman or Thunder Client to do tests in the meantime.



Notes on Thought Process
The RESTful API is designed to provide information about car brands and their models from a MongoDB database. The API uses the Gorilla Mux router to handle incoming requests and the MongoDB driver to communicate with the database.

The main.go file is responsible for starting the server and defining the server settings. The controller.go file contains the functions that handle incoming requests and interact with the database. The models.go file contains the structures that define the data stored in the database.

To test the API, the standard Go testing package is used. The tests are designed to cover all the functions in the controller.go file, and they ensure that the API is functioning correctly and returning the expected data.

Issues Encountered
There were no issues encountered while running, or testing the API. However, if you encounter any issues while using the API, please feel free to contact us for assistance.

The Issues encountered while building, where not coding issues but understanding the Coding Exercise arquitecture, the arquitecture given was not complete and some assumptions had to be made, 