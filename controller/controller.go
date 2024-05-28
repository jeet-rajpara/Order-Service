package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	amqp "order_service/amqp_helper"
	"order_service/model/request"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder request.Order

	body, err := io.ReadAll(r.Body)
	if err != nil {
		// errorhandling.SendErrorResponse(r, w, errorhandling.ReadBodyError, constant.EMPTY_STRING)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &newOrder)
	if err != nil {
		// errorhandling.HandleJSONUnmarshlError(r, w, err)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	ctx := r.Context()
	db := ctx.Value("db").(*sql.DB)

	var orderID int64
	query := "INSERT INTO orders (id,cus_name,cus_email,items,status) VALUES (DEFAULT,$1, $2, $3,$4) RETURNING id;"
	row := db.QueryRow(query, newOrder.CustomerName, newOrder.CustomerEmail, newOrder.Items, newOrder.Status)
	err = row.Scan(&orderID)
	if err != nil {
		log.Printf("Error inserting project: %v", err)
		return
	}
	fmt.Println(orderID)
	fmt.Println(newOrder.CustomerName) // Start a new Go routine to publish the order details to RabbitMQ

	go func(order request.Order) {
		conn, ch, err := amqp.ConnectRabbitMQ()
		if err != nil {
			log.Printf("Error connecting to RabbitMQ: %v", err)
			return
		}
		defer conn.Close()
		defer ch.Close()

		err = amqp.PublishMessage(ch, "notificationExchange", order)
		if err != nil {
			log.Printf("Error publishing to RabbitMQ: %v", err)
		}
	}(newOrder)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Order created with ID: %d", orderID)
}

