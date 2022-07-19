package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	nrgin "github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func makeGinEndpoint(s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Writer.WriteString(s)
	}
}

func v1login(c *gin.Context)  { c.Writer.WriteString("v1 login") }
func v1submit(c *gin.Context) { c.Writer.WriteString("v1 submit") }
func v1read(c *gin.Context)   { c.Writer.WriteString("v1 read") }

func endpoint404(c *gin.Context) {
	c.Writer.WriteHeader(404)
	c.Writer.WriteString("returning 404")
}

func endpointChangeCode(c *gin.Context) {
	// gin.ResponseWriter buffers the response code so that it can be
	// changed before the first write.
	c.Writer.WriteHeader(404)
	c.Writer.WriteHeader(200)
	c.Writer.WriteString("actually ok!")
}

func endpointResponseHeaders(c *gin.Context) {
	// Since gin.ResponseWriter buffers the response code, response headers
	// can be set afterwards.
	c.Writer.WriteHeader(200)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteString(`{"zip":"zap"}`)
}

func endpointNotFound(c *gin.Context) {
	c.Writer.WriteString("there's no endpoint for that!")
}

func endpointAccessTransaction(c *gin.Context) {
	txn := nrgin.Transaction(c)
	txn.SetName("custom-name")
	c.Writer.WriteString("changed the name of the transaction!")
}

func main() {

	apiKey, ok := os.LookupEnv("NEW_RELIC_API_KEY")
	if !ok {
		fmt.Println("Missing NEW_RELIC_API_KEY required for New Relic OpenTelemetry Exporter")
		os.Exit(1)
	}

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("newrelic-opentelemetry-poc GinServer"),
		newrelic.ConfigLicense(apiKey),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(nrgin.Middleware(app))

	router.GET("/404", endpoint404)
	router.GET("/change", endpointChangeCode)
	router.GET("/headers", endpointResponseHeaders)
	router.GET("/txn", endpointAccessTransaction)

	router.GET("/anon", func(c *gin.Context) {
		c.Writer.WriteString("anonymous function handler")
	})

	v1 := router.Group("/v1")
	v1.GET("/login", v1login)
	v1.GET("/submit", v1submit)
	v1.GET("/read", v1read)

	router.NoRoute(endpointNotFound)

	router.Run(":8000")
}
