package middleware

import (
	"bigdataconcept/fintech/intelligent/payment/routing/gateway/repository"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)



type RequestLogger struct {
	*repository.MongoRepo
}

func NewRequestLogger()  *RequestLogger{
	return &RequestLogger{}
}

func (RequestLogger *RequestLogger) LogIncomingRequest(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var bod = req.Body
		body, err := ioutil.ReadAll(bod)
		if err != nil {
			panic(err)
		}
		log.Println(string(body))
		req.Body.Close()  //
    		h.ServeHTTP(w, req)
	})
}