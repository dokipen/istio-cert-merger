package istiocertmerger

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Certificate struct {
	Name    string
	CrtData []byte
	KeyData []byte
	Hosts   []string
}

func (c *Certificate) CrtFilename() string {
	return fmt.Sprintf("%s.tls.crt", c.Name)
}

func (c *Certificate) KeyFilename() string {
	return fmt.Sprintf("%s.tls.key", c.Name)
}

func GetCerts(clientset *kubernetes.Clientset, namespace string, label string) ([]*Certificate, error) {
	var certs []*Certificate

	opts := metav1.ListOptions{LabelSelector: label}
	secrets, err := clientset.CoreV1().Secrets(namespace).List(opts)
	if err != nil {
		return certs, err
	}

	for _, secret := range secrets.Items {
		meta := secret.ObjectMeta
		certName := meta.Labels[CertificateNameTag]
		if certName == "" || secret.Data[TLSCrt] == nil || secret.Data[TLSKey] == nil {
			continue
		}

		certs = append(certs, &Certificate{
			Name:    certName,
			CrtData: secret.Data[TLSCrt],
			KeyData: secret.Data[TLSKey],
			Hosts:   strings.Split(meta.Annotations[AltNamesTag], ","),
		})
	}
	return certs, nil
}
