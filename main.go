package main

import (
	"context"
	"log"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()
	cleanUp := initTracer()
	defer cleanUp(ctx)

	router := gin.New()
	router.Use(otelgin.Middleware(serviceName))

	router.GET("/roll", func(c *gin.Context) {
		tracer := otel.Tracer("get roll")
		ctx, span := tracer.Start(c.Request.Context(), "rand.IntN")
		result := rand.IntN(20) + 1
		span.SetAttributes(attribute.KeyValue{Key: "result", Value: attribute.IntValue(result)})
		span.End()

		_, span = tracer.Start(ctx, "Serve JSON")
		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
		span.End()
	})

	log.Fatalln(router.Run(os.Getenv("GINOTEL_PORT")))
}
