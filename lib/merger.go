package istiocertmerger

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dokipen/istio-cert-merger/lib/spec"
)

type CertMerger struct {
	Certificates []*Certificate
	Spec         *spec.Spec
}

func NewCertMerger(certificates []*Certificate, spec *spec.Spec) *CertMerger {
	return &CertMerger{certificates, spec}
}

func (cm *CertMerger) getSecretData() map[string]string {
	data := map[string]string{}

	for _, cert := range cm.Certificates {
		data[cert.KeyFilename()] = base64.StdEncoding.EncodeToString(cert.KeyData)
		data[cert.CrtFilename()] = base64.StdEncoding.EncodeToString(cert.CrtData)
	}
	return data
}

func (cm *CertMerger) Revision() string {
	hash := sha256.New()
	for _, cert := range cm.Certificates {
		hash.Write(cert.KeyData)
		hash.Write(cert.CrtData)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (cm *CertMerger) GetPayload() ([]byte, error) {
	secretName := fmt.Sprintf("istio-%s-certs", cm.Spec.IstioIngress.Name)
	ingress := cm.Spec.IstioIngress
	response := map[string]interface{}{
		"attachments": []map[string]interface{}{
			{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"name":      secretName,
					"namespace": cm.Spec.TargetSecretNamespace,
					"annotations": map[string]string{
						RevisionTag: cm.Revision(),
					},
				},
				"data": cm.getSecretData(),
			},
			{
				"apiVersion": "networking.istio.io/v1alpha3",
				"kind":       "Gateway",
				"metadata": map[string]string{
					"name":      fmt.Sprintf("%s-cert-merge", ingress.Name),
					"namespace": ingress.Namespace,
				},
				"spec": map[string]interface{}{
					"selector": map[string]string{
						"istio": ingress.Name,
					},
					"servers": cm.getGatewayServers(),
				},
			},
		},
	}

	return json.Marshal(response)
}

func (cm *CertMerger) getGatewayServers() []*GatewayServer {
	var gs []*GatewayServer
	for _, cert := range cm.Certificates {
		gs = append(gs, cm.newGatewayServer(cert))
	}
	// Redirect server
	gs = append(gs, &GatewayServer{
		Hosts: []string{"*"},
		Port: &Port{
			Name:     "http",
			Number:   80,
			Protocol: "HTTP",
		},
		TLS: &TLS{
			HttpsRedirect: true,
		},
	})
	return gs
}

func (cm *CertMerger) newGatewayServer(cert *Certificate) *GatewayServer {
	ingress := cm.Spec.IstioIngress
	return &GatewayServer{
		Hosts: cert.Hosts,
		Port: &Port{
			Name:     fmt.Sprintf("%s-https", cert.Name),
			Number:   ingress.Port.Number,
			Protocol: ingress.Port.Protocol,
		},
		TLS: &TLS{
			Mode:              "SIMPLE",
			PrivateKey:        fmt.Sprintf("/etc/istio/%s-certs/%s.tls.key", ingress.Name, cert.Name),
			ServerCertificate: fmt.Sprintf("/etc/istio/%s-certs/%s.tls.crt", ingress.Name, cert.Name),
		},
	}
}
