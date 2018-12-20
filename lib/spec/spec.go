package spec

import (
	"encoding/json"
)

const (
	CertificateName = "certmanager.k8s.io/certificate-name"
)

func defaultSpec() *Spec {
	return &Spec{
		TargetSecretNamespace: "istio-system",
		Certificates: &Certificates{
			Namespace:     "cert-manager",
			LabelSelector: CertificateName,
		},
		IstioIngress: &IstioIngress{
			Namespace: "istio-system",
			Name:      "",
			Port: &Port{
				Name:     "",
				Number:   443,
				Protocol: "HTTPS",
			},
			SyncDeployment: true,
		},
	}
}

type Certificates struct {
	Namespace     string `json:"namespace"`
	LabelSelector string `json:"labelSelector"`
}

type Port struct {
	Name     string `json:"name"`
	Number   int32  `json:"number"`
	Protocol string `json:"protocol"`
}

type IstioIngress struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	Port           *Port  `json:"port"`
	SyncDeployment bool   `json:"syncDeployment"`
}

type Spec struct {
	TargetSecretNamespace string        `json:"targetSecretNamespace"`
	Certificates          *Certificates `json:"certificates"`
	IstioIngress          *IstioIngress `json:"istioIngress"`
}

type Request struct {
	Object *Object `json:"object"`
}

type Object struct {
	Spec *Spec `json:"spec"`
}

func ParseCertMergeRequest(body []byte) (*Spec, error) {
	request := Request{
		Object: &Object{
			Spec: defaultSpec(),
		},
	}
	err := json.Unmarshal(body, &request)
	return request.Object.Spec, err
}
