package proxy

import (
	"io"
	"log"
	"main/pkg/model/convert"
	"main/pkg/repository"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	CONNECT = "CONNECT"
	HTTP    = "http"
	HTTPS   = "https"
)

type Proxy struct {
	Crt, Key, Protocol string
	Repo               *repository.Repository
}

func (p *Proxy) StartProxy(server *http.Server) {
	switch p.Protocol {
	case HTTP:
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf(err.Error())
		}
		break
	case HTTPS:
		if err := server.ListenAndServeTLS(p.Crt, p.Key); err != nil {
			log.Fatalf(err.Error())
		}
		break
	default:
		log.Println("not http or https")
		break
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == CONNECT {
		p.handleHTTPS(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}

func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	request := convert.ParseHTTPRequest(r)
	response := &http.Response{}

	r.RequestURI = ""
	r.Header.Del("Proxy-Connection")

	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	proxyResponse, err := httpClient.Do(r)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer proxyResponse.Body.Close()

	response.Header = make(http.Header)
	for header, values := range proxyResponse.Header {
		stringValues := strings.Join(values, ", ")
		w.Header().Set(header, stringValues)
		response.Header.Set(header, stringValues)
	}
	w.WriteHeader(proxyResponse.StatusCode)
	response.StatusCode = proxyResponse.StatusCode

	response.Body = proxyResponse.Body
	io.Copy(w, proxyResponse.Body)

	request.Response = *convert.ParseHTTPResponse(response)
	_, err = p.Repo.RequestRepository.CreateRequest(r.Context(), request)
	if err != nil {
		log.Printf("Don`t save request with ID: %s", request.ID)
	}
}

func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	//request := convert.ParseHTTPRequest(r)

	connDest, err := net.DialTimeout("tcp", r.Host, time.Second)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connSrc, _, err := hijacker.Hijack()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go broadcastData(connDest, connSrc)
	go broadcastData(connSrc, connDest)
}

func broadcastData(to io.WriteCloser, from io.ReadCloser) {
	io.Copy(to, from)
	to.Close()
	from.Close()
}
