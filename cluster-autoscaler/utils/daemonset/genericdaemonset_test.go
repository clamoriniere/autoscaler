package daemonset

import (
	"reflect"
	"strings"
	"testing"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenericDaemonSet_GetPodTemplateSpec(t *testing.T) {
	defaultPodTemplate := &apiv1.PodTemplateSpec{
		Spec: apiv1.PodSpec{
			NodeName:           "nodename-foo",
			ServiceAccountName: "service-account-foo",
		},
	}

	podTemplateMap := structToMap(defaultPodTemplate)
	specMap := map[string]interface{}{
		"template": podTemplateMap,
	}

	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		argPath string
		want    *apiv1.PodTemplateSpec
		wantErr bool
	}{
		{
			name:    "wrong podTemplate path",
			argPath: "foo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "decode OK, path ok",
			argPath: "template",
			fields: fields{
				Spec: specMap,
			},
			want:    defaultPodTemplate.DeepCopy(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &GenericDaemonSet{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
			}
			got, err := d.GetPodTemplateSpec(tt.argPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenericDaemonSet.GetPodTemplateSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenericDaemonSet.GetPodTemplateSpec() = %v,\n want %v", got, tt.want)
			}
		})
	}
}

func structToMap(item interface{}) map[string]interface{} {

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		tag = strings.Split(tag, ",")[0] // remove omitempty from tag
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
