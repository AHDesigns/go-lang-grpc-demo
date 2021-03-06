package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/user/learn-proto/proto"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure());
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	g := gin.Default()

	g.GET("/add/:a/:b", func(ctx *gin.Context){
		a := getParamOfFail(ctx, "a")
		b := getParamOfFail(ctx, "b")

		req := &proto.Request{A: int64(a), B: int64(b)}
		if response, err := client.Add(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"result": fmt.Sprint(response.Result)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	g.GET("/mult/:a/:b", func(ctx *gin.Context){
		a := getParamOfFail(ctx, "a")
		b := getParamOfFail(ctx, "b")

		req := &proto.Request{A: int64(a), B: int64(b)}

		if response, err := client.Multiply(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{"request": fmt.Sprint(response.Result)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run Server: %v", err)
	}
}

func getParamOfFail(ctx *gin.Context, key string) uint64 {
	param, err := strconv.ParseUint(ctx.Param(key), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid Parameter %s", key)})
	}
	return param
}
