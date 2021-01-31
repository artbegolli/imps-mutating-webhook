package mutators

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

type ServiceAccountMutator struct {
	Log     *logrus.Logger
	KubeCli kubernetes.Interface
}

func (s *ServiceAccountMutator) Mutate(review *admissionv1.AdmissionReview) (*admissionv1.AdmissionReview, error) {

	ctx := context.Background()
	cm, err := s.KubeCli.CoreV1().ConfigMaps("kube-system").Get(ctx, "imps-mutating-webhook-configmap", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pullSecretName := cm.Data["pull-secret-ref"]

	serviceAccount := &corev1.ServiceAccount{}
	if err := json.Unmarshal(review.Request.Object.Raw, serviceAccount); err != nil {
		return nil, fmt.Errorf("unable unmarshal service account json object %v", err)
	}

	patch := map[string]string{
		"op":    "add",
		"path":  "/imagePullSecrets",
		"value": fmt.Sprintf("[\"name\": \"%s\"]", pullSecretName),
	}
	review.Response.Patch, err = json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	review.Response.Result = &metav1.Status{
		Status: "Success",
	}
	review.Response.Allowed = true
	review.Response.UID = review.Request.UID

	jsonPatchType := admissionv1.PatchTypeJSONPatch
	review.Response.PatchType = &jsonPatchType

	return review, nil
}
