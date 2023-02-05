package main
import (
    "fmt"
    "net/http"
    "io"
    "time"
    "log"
    //"net/url"
    // "net"
    "strings"
    "strconv"
)

const (
        proxyPort   = 11095
        servicePort = 80
)

type Proxy struct{}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Forward the HTTP request to the destination service.
        res, duration, err := p.forwardRequest(req)

	// Notify the client if there was an error while forwarding the request.
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadGateway)
                return
        }

	// If the request was forwarded successfully, write the response back to
	// the client.
        p.writeResponse(w, res)

	// Print request and response statistics.
        p.printStats(req, res, duration)
}

func (p *Proxy) forwardRequest(req *http.Request) (*http.Response, time.Duration, error) {

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
	log.Printf("\nIncoming URL: %s\nForward URL: %s\nMethod: %s", incUrl, intUrl, req.Method)

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

func (p *Proxy) writeResponse(w http.ResponseWriter, res *http.Response) {
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

func (p *Proxy) printStats(req *http.Request, res *http.Response, duration time.Duration) {
        fmt.Printf("Request Duration: %v\n", duration)
        fmt.Printf("Request Size: %d\n", req.ContentLength)
        fmt.Printf("Response Size: %d\n", res.ContentLength)
        fmt.Printf("Response Status: %d\n\n", res.StatusCode)
}

func main() {
	log.Printf("Starting forwarding proxy. Listening at port: %d; Sending to port:%d", proxyPort, servicePort)

	http.ListenAndServe(fmt.Sprintf(":%d", proxyPort), &Proxy{})
}

