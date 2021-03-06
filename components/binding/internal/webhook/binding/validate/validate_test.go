package validate

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/kyma-project/kyma/components/binding/pkg/apis/v1alpha1"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	namespace = "test"
)

func TestValidationHandler_Handle(t *testing.T) {
	// given
	sch, err := v1alpha1.SchemeBuilder.Build()
	require.NoError(t, err)
	err = scheme.AddToScheme(sch)
	require.NoError(t, err)

	binding := fixBinding("binding-test")
	rawBinding, err := json.Marshal(binding)
	require.NoError(t, err)

	request := admission.Request{
		AdmissionRequest: v1beta1.AdmissionRequest{
			UID:       "1234-abcd",
			Operation: v1beta1.Create,
			Name:      "test-binding",
			Namespace: namespace,
			Kind: metav1.GroupVersionKind{
				Kind:    "Binding",
				Version: "v1alpha1",
				Group:   "bindings.kyma-project.io",
			},
			Object: runtime.RawExtension{Raw: rawBinding},
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(sch)
	decoder, err := admission.NewDecoder(scheme.Scheme)
	require.NoError(t, err)

	handler := NewValidationHandler(logrus.New())
	err = handler.InjectClient(fakeClient)
	require.NoError(t, err)
	err = handler.InjectDecoder(decoder)
	require.NoError(t, err)

	// when
	response := handler.Handle(context.TODO(), request)

	// then
	assert.True(t, response.Allowed)
}

func TestValidationHandler_HandleError(t *testing.T) {
	// given
	sch, err := v1alpha1.SchemeBuilder.Build()
	require.NoError(t, err)
	err = scheme.AddToScheme(sch)
	require.NoError(t, err)

	binding := fixBinding("binding_test")
	rawBinding, err := json.Marshal(binding)
	require.NoError(t, err)

	request := admission.Request{
		AdmissionRequest: v1beta1.AdmissionRequest{
			UID:       "1234-abcd",
			Operation: v1beta1.Create,
			Name:      "test-binding",
			Namespace: namespace,
			Kind: metav1.GroupVersionKind{
				Kind:    "Binding",
				Version: "v1alpha1",
				Group:   "bindings.kyma-project.io",
			},
			Object: runtime.RawExtension{Raw: rawBinding},
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(sch)
	decoder, err := admission.NewDecoder(scheme.Scheme)
	require.NoError(t, err)

	handler := NewValidationHandler(logrus.New())
	err = handler.InjectClient(fakeClient)
	require.NoError(t, err)
	err = handler.InjectDecoder(decoder)
	require.NoError(t, err)

	// when
	response := handler.Handle(context.TODO(), request)

	// then
	assert.False(t, response.Allowed)
}

func fixBinding(name string) *v1alpha1.Binding {
	return &v1alpha1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
}
