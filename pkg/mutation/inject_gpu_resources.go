package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// injectGPUResource is a container for the mutation injecting environment vars
type injectGPUResource struct {
	Logger logrus.FieldLogger
}

// injectGPUResource implements the podMutator interface
var _ podMutator = (*injectGPUResource)(nil)

// Name returns the struct name
func (se injectGPUResource) Name() string {
	return "inject_env"
}

// Mutate returns a new mutated pod according to set env rules
func (se injectGPUResource) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation - injecting GPU resources", se.Name())
	mpod := pod.DeepCopy()
	// se.Logger.Debugf("Adding resources %s", mgpu)
	injectGPUResourceToContainer(mpod)
	return mpod, nil
}

// injectGPUResourceToContainer assigns gpu to build container
func injectGPUResourceToContainer(pod *corev1.Pod) {
	for i, container := range pod.Spec.Containers {
		if container.Name == "build" {
			// Ensure the current requests map is initialized
			if pod.Spec.Containers[i].Resources.Requests == nil {
				pod.Spec.Containers[i].Resources.Requests = corev1.ResourceList{}
			}
			if pod.Spec.Containers[i].Resources.Limits == nil {
				pod.Spec.Containers[i].Resources.Limits = corev1.ResourceList{}
			}

			// Add GPU resource request
			gpu := resource.MustParse("1")
			if _, ok := pod.Spec.Containers[i].Resources.Requests["nvidia.com/gpu"]; !ok {
				pod.Spec.Containers[i].Resources.Requests["nvidia.com/gpu"] = gpu
				pod.Spec.Containers[i].Resources.Limits["nvidia.com/gpu"] = gpu
			}
		}
	}
}
