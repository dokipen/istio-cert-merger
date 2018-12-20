package istiocertmerger

import (
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/dokipen/istio-cert-merger/lib/spec"
)

type Server struct {
	Clientset *kubernetes.Clientset
}

// TODO(bob): Path and Content-Type should be validated.
func (s *Server) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	logrus.WithFields(logrus.Fields{"body": string(body)}).Debug("Request recieved")
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	spec, err := spec.ParseCertMergeRequest(body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	if spec.IstioIngress.Name == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("object.spec.istioIngress.name must be defined in JSON payload\n"))
		return
	}
	certs, err := GetCerts(s.Clientset, spec.Certificates.Namespace, spec.Certificates.LabelSelector)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("Failed to get certs")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	merger := NewCertMerger(certs, spec)

	payload, err := merger.GetPayload()
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Error("Failed to merge certs")
		responseWriter.WriteHeader(http.StatusBadRequest)
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(payload)

	ingress := spec.IstioIngress
	if ingress.SyncDeployment {
		err = SyncIngressDeployment(spec, s.Clientset)
		if err != nil {
			logrus.WithFields(logrus.Fields{"err": err}).Error("Failed to sync ingress deployment")
		} else {
			logrus.Info("Ingress deployment sync completed")
		}
	}
}
