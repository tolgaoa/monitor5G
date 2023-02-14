package main

import (
	"context"
	"fmt"
	"log"
	//"os"
	//"os/signal"
	"time"
	"net/http"
	"io"
	"os"
	"strings"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
        proxyPort   = 11095
        servicePort = 80
)

type Proxy struct{}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initTracer() (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(os.Getenv("SERVICENAME")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	//conn, err := grpc.DialContext(ctx, "10.42.0.14:4317",
	conn, err := grpc.DialContext(ctx, "otel-collector-daemonset-collector.otel-collector.svc.cluster.local:4317",
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

func forwardRequest(req *http.Request) (*http.Response, time.Duration, error) {

        // Prepare the destination endpoint to forward the request to.
        //proxyUrl := fmt.Sprintf("http://%s:%d%s", req.Host, servicePort, req.RequestURI)
        incUrl := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)

        /*
        u, _ := url.Parse(incUrl)
        fmt.Println("full uri:", u.String())
        fmt.Println("scheme:", u.Scheme)
        fmt.Println("opaque:", u.Opaque)
        fmt.Println("Host:", u.Host)
        fmt.Println("Path", u.Path)
        fmt.Println("Fragment", u.Fragment)
        fmt.Println("RawQuery", u.RawQuery)
        fmt.Printf("query: %#v", u.Query())
        fmt.Printf("\n")
        thost, tport, _ := net.SplitHostPort(req.Host)
        urlsplit := strings.Split(u.String(), ":")
        */


        intUrl := strings.Replace(incUrl, strconv.Itoa(proxyPort), strconv.Itoa(servicePort), 1)

        // Print the original URL and the proxied request URL.
        bodyBytes, err := io.ReadAll(req.Body)
        bodyString := string(bodyBytes)
        log.Printf("\nIncoming URL: %s\nForward URL: %s\nMethod: %s\nBody: %s", incUrl, intUrl, req.Method, bodyString)

        // Create an HTTP client and a proxy request based on the original request.
        httpClient := http.Client{}
        proxyReq, err := http.NewRequest(req.Method, intUrl, req.Body)

        // Capture the duration while making a request to the destination service.
        start := time.Now()
        res, err := httpClient.Do(proxyReq)
        duration := time.Since(start)

        // Return the response, the request duration, and the error.
        return res, duration, err
}

func writeResponse(w http.ResponseWriter, res *http.Response) {
        // Copy all the header values from the response.
        for name, values := range res.Header {
                w.Header()[name] = values
        }

        // Set a special header to notify that the proxy actually serviced the request.
        w.Header().Set("Server", "amazing-proxy")

        // Set the status code returned by the destination service.
        w.WriteHeader(res.StatusCode)

        // Copy the contents from the response body.
        io.Copy(w, res.Body)

        // Finish the request.
        res.Body.Close()
}

func printStats(req *http.Request, res *http.Response, duration time.Duration) {
        fmt.Printf("Request Duration: %v\n", duration)
        fmt.Printf("Request Size: %d\n", req.ContentLength)
        fmt.Printf("Response Size: %d\n", res.ContentLength)
        fmt.Printf("Response Status: %d\n\n", res.StatusCode)
}


func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	uk := attribute.Key("username")

	defaultHandler := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		span := trace.SpanFromContext(ctx)
		bag := baggage.FromContext(ctx)
		span.AddEvent("handling this...", trace.WithAttributes(uk.String(bag.Member("username").Value())))

		// Forward the HTTP request to the destination service.
		res, duration, err := forwardRequest(req)

		// Notify the client if there was an error while forwarding the request.
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// If the request was forwarded successfully, write the response back to
		// the client.
		writeResponse(w, res)

		// Print request and response statistics.
		printStats(req, res, duration)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(defaultHandler), "Default")

	http.Handle("/", otelHandler)
	log.Printf("Starting HTTP server")
	err = http.ListenAndServe(":11095", nil)
	if err != nil {
		log.Fatal(err)
	}
}
