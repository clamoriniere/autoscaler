package daemonset

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenericDaemonSet represent a generic Daemonset Resource
type GenericDaemonSet struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// The desired behavior of this daemon set.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec map[string]interface{} `json:"spec,omitempty"`
}

// GetPodTemplateSpec used to access the PodTemplateSpec from the GenericDaemonSetSpec
func (d *GenericDaemonSet) GetPodTemplateSpec(path string) (*apiv1.PodTemplateSpec, error) {
	if d.Spec == nil {
		return nil, fmt.Errorf("Spec is nil")
	}

	jsonPath := strings.Split(path, ".")

	var podTemplateSpecMap map[string]interface{}
	podTemplateSpecMap = d.Spec
	for _, path := range jsonPath {
		tmpPath, ok := podTemplateSpecMap[path]
		if !ok {
			return nil, fmt.Errorf("unable to access the podTemplate path, %s", path)
		}
		podTemplateSpecMap = tmpPath.(map[string]interface{})
	}

	var podTemplateSpec apiv1.PodTemplateSpec
	cfg := &mapstructure.DecoderConfig{
		Result:  &podTemplateSpec,
		TagName: "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err := decoder.Decode(podTemplateSpecMap)

	return &podTemplateSpec, err
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GenericDaemonSetList is a collection of daemon sets.
type GenericDaemonSetList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// A list of daemon sets.
	Items []GenericDaemonSet `json:"items" protobuf:"bytes,2,rep,name=items"`
}
