package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/artbegolli/imps-mutating-webhook/pkg/log"
	"github.com/artbegolli/imps-mutating-webhook/pkg/mutators"
	v1 "k8s.io/api/admission/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func HandleServiceAccountMutate(w http.ResponseWriter, req *http.Request) {

	logger := log.GetLogger(true)

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.WithError(err).Errorln("error creating in cluster config")
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.WithError(err).Errorln("error creating in kubernetes client")
		return
	}

	mutator := mutators.ServiceAccountMutator{
		Log:     logger,
		KubeCli: clientset,
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.WithError(err).Errorln("error reading request body")
		return
	}

	review := &v1.AdmissionReview{}
	if err := json.Unmarshal(body, review); err != nil {
		logger.WithError(err).Errorln("error unmarshalling admission review")
		return
	}

	mutatedAdmissionReview, err := mutator.Mutate(review)
	if err != nil {
		logger.WithError(err).Errorln("error mutating resource")
		return
	}

	mutatedBody, err := json.Marshal(mutatedAdmissionReview)
	if err != nil {
		logger.WithError(err).Errorln("error marshalling mutated body")
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write(mutatedBody); err != nil {
		logger.WithError(err).Errorln("error writing mutated body to response")
		return
	}
}

func HandlePing(w http.ResponseWriter, req *http.Request) {
	fmt.Println("ping")
}

func main() {

	http.HandleFunc("/mutate", HandleServiceAccountMutate)
	http.HandleFunc("/ping", HandlePing)
	fmt.Println("Listening for requests on :443")
	if err := http.ListenAndServeTLS(":443", "/crts/tls.crt", "/crts/tls.key", nil); err != nil {
		panic("error starting TLS server" + err.Error())
		return
	}

}
