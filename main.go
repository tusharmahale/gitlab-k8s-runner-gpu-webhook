package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tusharmahale/container-admission-webhook/pkg/admission"
	admissionv1 "k8s.io/api/admission/v1"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mutate-webhook", ServeMutatePods).Methods(http.MethodPost)
	router.HandleFunc("/healthz", healthCheck).Methods(http.MethodGet)
	router.NotFoundHandler = http.HandlerFunc(invalidURL)
	if os.Getenv("TLS") == "true" {
		cert := "/etc/webhook/tls/tls.crt"
		key := "/etc/webhook/tls/tls.key"
		logrus.Print("Listening on port 443...")
		err := http.ListenAndServeTLS(":443", cert, key, router)
		if err != nil {
			log.Fatal("Error starting server ", err)
		}
	} else {
		logrus.Print("Listening on port 8080...")
		err := http.ListenAndServe(":8080", router)
		if err != nil {
			log.Fatal("Error starting server ", err)
		}
	}

}

func invalidURL(w http.ResponseWriter, r *http.Request) {
	returnPl := map[string]string{}
	w.Header().Set("Content-type", "application/json")
	returnPl["error"] = "invalid URL"
	w.WriteHeader(http.StatusNotFound)
	returnBytes, _ := json.Marshal(returnPl)
	w.Write(returnBytes)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	returnPl := map[string]string{}
	w.Header().Set("Content-type", "application/json")
	returnPl["status"] = "success"
	w.WriteHeader(http.StatusOK)
	returnBytes, _ := json.Marshal(returnPl)
	w.Write(returnBytes)
}

// ServeMutatePods returns an admission review with pod mutations as a json patch
// in the review response
func ServeMutatePods(w http.ResponseWriter, r *http.Request) {
	logger := logrus.WithField("uri", r.RequestURI)
	logger.Debug("received mutation request")

	in, err := parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adm := admission.Admitter{
		Logger:  logger,
		Request: in.Request,
	}

	out, err := adm.MutatePodReview()
	if err != nil {
		e := fmt.Sprintf("could not generate admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		e := fmt.Sprintf("could not parse admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	logger.Debug("sending response")
	logger.Debugf("%s", jout)
	fmt.Fprintf(w, "%s", jout)
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}
	fmt.Println("Request parsed successfully")

	return &a, nil
}
