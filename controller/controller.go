package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlirGG/mongoapi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

const connectionString = "mongodb+srv://user_test:password_test@alirggtestscluster.hlxprlo.mongodb.net/?retryWrites=true&w=majority"

func connectToDB() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetBrands(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	brandsCollection := client.Database("nexuBrands").Collection("brands")
	cursor, err := brandsCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Create a map to store the sum of the model average prices for each brand
	brandSums := make(map[string]float64)
	// Create a map to store the number of models for each brand
	brandCounts := make(map[string]int)

	for cursor.Next(ctx) {
		var model model.Brand
		if err := cursor.Decode(&model); err != nil {
			log.Fatal(err)
		}

		// Update the sum and count for the brand
		brandSums[model.Name] += model.AveragePrice
		brandCounts[model.Name]++
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Create a slice to store the output brands
	var brands []map[string]interface{}

	// Loop through each brand in the sums map and create an output brand object
	for brand, sum := range brandSums {
		count := brandCounts[brand]
		averagePrice := sum / float64(count)
		outputBrand := map[string]interface{}{
			"id":            len(brands) + 1,
			"nombre":        brand,
			"average_price": averagePrice,
		}
		brands = append(brands, outputBrand)
	}

	json.NewEncoder(w).Encode(brands)
}

func GetModelsForBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	brandsCollection := client.Database("nexuBrands").Collection("brands")

	params := strings.Split(r.URL.Path, "/")
	brandName := params[len(params)-2]

	filter := bson.M{"brand_name": brandName}

	cur, err := brandsCollection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Brand not found")
		return
	}

	var models []map[string]interface{}
	for cur.Next(ctx) {
		var model struct {
			ID           int    `bson:"id"`
			Name         string `bson:"name"`
			AveragePrice int    `bson:"average_price"`
		}
		if err := cur.Decode(&model); err != nil {
			log.Fatal(err)
		}

		modelMap := map[string]interface{}{
			"id":            model.ID,
			"name":          model.Name,
			"average_price": model.AveragePrice,
		}

		models = append(models, modelMap)
	}

	if len(models) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Brand not found")
		return
	}

	json.NewEncoder(w).Encode(models)
}

func CreateBrand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get brand data from request body
	var brand model.Brand
	err := json.NewDecoder(r.Body).Decode(&brand)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Check if brand name already exists
	client, err := connectToDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error connecting to database"})
		return
	}
	defer client.Disconnect(ctx)

	brandsCollection := client.Database("nexuBrands").Collection("brands")

	filter := bson.M{"brand_name": brand.Name}

	var existingBrand model.Brand
	err = brandsCollection.FindOne(ctx, filter).Decode(&existingBrand)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "Brand already exists"})
		return
	}

	// Generate new brand ID
	cursor, err := brandsCollection.Find(ctx, bson.D{{}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error generating brand ID"})
		return
	}
	defer cursor.Close(ctx)

	var brands []model.Brand
	if err = cursor.All(ctx, &brands); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error generating brand ID"})
		return
	}

	var maxID int64
	for _, brand := range brands {
		if brand.ID > maxID {
			maxID = brand.ID
		}
	}
	brand.ID = maxID + 1

	// Add new brand to database
	_, err = brandsCollection.InsertOne(ctx, brand)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error adding brand to database"})
		return
	}

	// Send response with new brand
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(brand)
}

func CreateModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get brand name and model name from URL parameters
	vars := mux.Vars(r)
	brandName := vars["brand_name"]
	modelName := vars["name"]

	// Get model data from request body
	var newModel model.Model
	err := json.NewDecoder(r.Body).Decode(&newModel)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Check if average price is given and greater than 100,000
	if newModel.AveragePrice != 0 && newModel.AveragePrice < 100000 {
		newModel.AveragePrice = 0
	}

	// Check if brand exists
	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	brandsCollection := client.Database("nexuBrands").Collection("brands")

	filter := bson.M{"brand_name": brandName}

	var existingBrand model.Brand
	err = brandsCollection.FindOne(ctx, filter).Decode(&existingBrand)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Brand not found"})
		return
	}

	// Check if model name already exists for this brand
	for _, m := range existingBrand.Models {
		if m.Name == modelName {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Model already exists"})
			return
		}
	}

	// Add new model to database
	newModel.ID = time.Now().UnixNano()
	existingBrand.Models = append(existingBrand.Models, newModel)
	update := bson.M{"$set": bson.M{"models": existingBrand.Models}}
	_, err = brandsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error adding model to database"})
		return
	}

	// Send response with new model
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newModel)
}

// PROBLEM WITH MONGODB ids no time to do it this way

func UpdateModel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get ID from URL parameters
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
		return
	}

	// Get new average price from request body
	var newPrice struct {
		AveragePrice float64 `json:"average_price"`
	}
	err = json.NewDecoder(r.Body).Decode(&newPrice)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	// Check if new average price is valid
	if newPrice.AveragePrice < 100000 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "New average price is not valid"})
		return
	}

	// Check if model exists
	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	modelsCollection := client.Database("nexuBrands").Collection("brands")
	filter := bson.M{"id": id}

	var existingModel model.Model
	err = modelsCollection.FindOne(ctx, filter).Decode(&existingModel)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Model not found"})
		return
	}

	// Update model with new average price
	existingModel.AveragePrice = newPrice.AveragePrice
	update := bson.M{"$set": bson.M{"average_price": existingModel.AveragePrice}}
	_, err = modelsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error updating model in database"})
		return
	}

	// Send response with updated model
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]float64{"average_price": existingModel.AveragePrice})
}

func GetAllModels(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	greaterParam := r.URL.Query().Get("greater")
	lowerParam := r.URL.Query().Get("lower")

	// Connect to database
	client, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Get collection
	modelsCollection := client.Database("nexuBrands").Collection("brands")

	// Build filter based on query params
	filter := bson.M{}
	if greaterParam != "" && lowerParam != "" {
		greaterPrice, err := strconv.ParseFloat(greaterParam, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid greater parameter"})
			return
		}
		lowerPrice, err := strconv.ParseFloat(lowerParam, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid lower parameter"})
			return
		}
		filter["average_price"] = bson.M{"$gt": greaterPrice, "$lt": lowerPrice}
	} else {
		if greaterParam != "" {
			greaterPrice, err := strconv.ParseFloat(greaterParam, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid greater parameter"})
				return
			}
			filter["average_price"] = bson.M{"$gt": greaterPrice}
		}
		if lowerParam != "" {
			lowerPrice, err := strconv.ParseFloat(lowerParam, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid lower parameter"})
				return
			}
			if filter["average_price"] != nil {
				filter["average_price"].(bson.M)["$lt"] = lowerPrice
			} else {
				filter["average_price"] = bson.M{"$lt": lowerPrice}
			}
		}

	}
	// Query database
	cur, err := modelsCollection.Find(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error querying database"})
		return
	}
	defer cur.Close(ctx)

	// Parse results into slice of models
	var models []model.Model
	for cur.Next(ctx) {
		var m model.Model
		err := cur.Decode(&m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error decoding model"})
			return
		}
		models = append(models, m)
	}
	if err := cur.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error iterating over models"})
		return
	}

	// Send response with slice of models
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}
