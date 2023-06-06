package main

import (
        "context"
        "log"
        "net/http"

	"revproxy/pkg/tracer"
	"revproxy/pkg/servcl"

        "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
        "go.opentelemetry.io/otel/baggage"
        "go.opentelemetry.io/otel/attribute"
        "go.opentelemetry.io/otel/trace"




)

func main() {
	tp, err := tracer.InitTracer()
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
		res, duration, err := servcl.ForwardRequest(req)

		// Notify the client if there was an error while forwarding the request.
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// If the request was forwarded successfully, write the response back to
		// the client.
		servcl.WriteResponse(w, res)

		// Print request and response statistics.
		servcl.PrintStats(req, res, duration)
	}

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(defaultHandler), "Default")

	http.Handle("/", otelHandler)
	log.Printf("Starting HTTP server")
	err = http.ListenAndServe(":11095", nil)
	if err != nil {
		log.Fatal(err)
	}
}
