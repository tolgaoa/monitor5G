package servcl

import (
        "fmt"
        "log"
        "time"
        "net/http"
        "io"
        "strings"
        "strconv"
        "io/ioutil"
        "bytes"
)

type Proxy struct{}

const (
        proxyPort   = 11095
        servicePort = 8080
)

func ForwardRequest(req *http.Request) (*http.Response, time.Duration, error) {

        // Prepare the destination endpoint to forward the request to.
        incUrl := fmt.Sprintf("http://%s%s", req.Host, req.RequestURI)
        intUrl := strings.Replace(incUrl, strconv.Itoa(proxyPort), strconv.Itoa(servicePort), 1)

        // Print the original URL and the proxied request URL.
        bodyBytes, err := ioutil.ReadAll(req.Body)
        rdr1 := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
        bodyString := string(bodyBytes)
        log.Printf("\nIncoming URL: %s\nForward URL: %s\nMethod: %s\nBody: %s", incUrl, intUrl, req.Method, bodyString)

        // Create an HTTP client and a proxy request based on the original request.
        httpClient := http.Client{}
        proxyReq, err := http.NewRequest(req.Method, intUrl, rdr1)

        // Capture the duration while making a request to the destination service.
        start := time.Now()
        res, err := httpClient.Do(proxyReq)
        duration := time.Since(start)

        // Return the response, the request duration, and the error.
        return res, duration, err
}

func WriteResponse(w http.ResponseWriter, res *http.Response) {
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

func PrintStats(req *http.Request, res *http.Response, duration time.Duration) {
        fmt.Printf("Request Duration: %v\n", duration)
        fmt.Printf("Request Size: %d\n", req.ContentLength)
        fmt.Printf("Response Size: %d\n", res.ContentLength)
        fmt.Printf("Response Status: %d\n\n", res.StatusCode)
}
