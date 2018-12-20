package istiocertmerger

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/dokipen/istio-cert-merger/lib/spec"
)

type IngressDeployment struct {
	Spec      *spec.Spec
	Clientset *kubernetes.Clientset
}

func (id *IngressDeployment) Sync() error {
	ingress := id.Spec.IstioIngress

	secrets := id.Clientset.CoreV1().Secrets(ingress.Namespace)
	secret, err := secrets.Get(fmt.Sprintf("istio-%s-certs", ingress.Name), metav1.GetOptions{})
	if err != nil {
		return err
	}
	revision := secret.ObjectMeta.Annotations[RevisionTag]

	deployments := id.Clientset.AppsV1().Deployments(ingress.Namespace)
	deployment, err := deployments.Get(fmt.Sprintf("istio-%s", ingress.Name), metav1.GetOptions{})
	if err != nil {
		return err
	}

	deployedRevision := deployment.Spec.Template.ObjectMeta.Annotations[RevisionTag]
	if deployedRevision != revision {
		logrus.WithFields(logrus.Fields{
			"namespace": id.Spec.IstioIngress.Namespace,
			"name":      id.Spec.IstioIngress.Name,
		}).Info("Syncing ingress deployment")
		deployment.Spec.Template.ObjectMeta.Annotations[RevisionTag] = revision
		deployments.Update(deployment)
	}
	return nil
}

func SyncIngressDeployment(spec *spec.Spec, clientset *kubernetes.Clientset) error {
	id := &IngressDeployment{
		Spec:      spec,
		Clientset: clientset,
	}
	return id.Sync()
}
