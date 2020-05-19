package main

import (
	"fmt"
	"hello/proto"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Handler of the order endpoint. Processing REST calls and invokes OrderService using gRPC. Returns order summary in JSON.
func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewOrderServiceClient(conn)
	g := gin.Default()

	handleRequests(g, client)

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func handleRequests(g *gin.Engine, client proto.OrderServiceClient) {

	g.POST("/order", func(ctx *gin.Context) {

		var order proto.Order
		err := ctx.BindJSON(&order)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		order.Id = int64(os.Getpid())

		log.Println("Submitting order:", order.Id)

		if resp, err := client.SubmitOrder(ctx, &order); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"amount":   fmt.Sprint(resp.GetAmount()),
				"order_id": order.Id,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("order %d failed", order.Id)})
	})

}
