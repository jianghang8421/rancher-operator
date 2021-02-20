package unmanaged

import (
	"context"
	"encoding/json"
	"strings"

	rancherv1 "github.com/rancher/rancher-operator/pkg/apis/rancher.cattle.io/v1"
	rkev1 "github.com/rancher/rancher-operator/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/rancher-operator/pkg/clients"
	capicontrollers "github.com/rancher/rancher-operator/pkg/generated/controllers/cluster.x-k8s.io/v1alpha4"
	mgmtcontroller "github.com/rancher/rancher-operator/pkg/generated/controllers/management.cattle.io/v3"
	rkecontroller "github.com/rancher/rancher-operator/pkg/generated/controllers/rke.cattle.io/v1"
	"github.com/rancher/rancher-operator/pkg/planner"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/data"
	"github.com/rancher/wrangler/pkg/data/convert"
	corecontrollers "github.com/rancher/wrangler/pkg/generated/controllers/core/v1"
	"github.com/rancher/wrangler/pkg/kv"
	corev1 "k8s.io/api/core/v1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	capi "sigs.k8s.io/cluster-api/api/v1alpha4"
)

const (
	machineRequestType = "rke.cattle.io/machine-request"
)

func Register(ctx context.Context, clients *clients.Clients) {
	h := handler{
		unmanagedMachine: clients.RKE.UnmanagedMachine(),
		mgmtClusterCache: clients.Management.Cluster().Cache(),
		capiClusterCache: clients.CAPI.Cluster().Cache(),
		machineCache:     clients.CAPI.Machine().Cache(),
		secrets:          clients.Core.Secret(),
		apply: clients.Apply.WithSetID("unmanaged-machine").
			WithCacheTypes(
				clients.Management.Cluster(),
				clients.Cluster.Cluster(),
				clients.RKE.UnmanagedMachine(),
				clients.CAPI.Machine(),
				clients.RKE.RKEBootstrap()),
	}
	clients.RKE.UnmanagedMachine().OnChange(ctx, "unmanaged-machine", h.onUnmanagedMachineChange)
	clients.Core.Secret().OnChange(ctx, "unmanaged-machine", h.onSecretChange)
}

type handler struct {
	unmanagedMachine rkecontroller.UnmanagedMachineClient
	mgmtClusterCache mgmtcontroller.ClusterCache
	capiClusterCache capicontrollers.ClusterCache
	machineCache     capicontrollers.MachineCache
	secrets          corecontrollers.SecretClient
	apply            apply.Apply
}

func (h *handler) onSecretChange(key string, secret *corev1.Secret) (*corev1.Secret, error) {
	if secret == nil || secret.Type != machineRequestType {
		return secret, nil
	}

	data := data.Object{}
	if err := json.Unmarshal(secret.Data["data"], &data); err != nil {
		// ignore invalid json, wait until it's valid
		return secret, nil
	}

	capiCluster, err := h.getCAPICluster(secret)
	if apierror.IsNotFound(err) {
		return secret, nil
	} else if err != nil {
		return secret, err
	}

	_, err = h.machineCache.Get(capiCluster.Namespace, secret.Name)
	if apierror.IsNotFound(err) {
		err = h.createMachine(capiCluster, secret, data)
	}
	if err != nil {
		return secret, err
	}

	if secret.Labels[planner.MachineNamespaceLabel] != capiCluster.Namespace ||
		secret.Labels[planner.MachineNameLabel] != secret.Name {
		secret = secret.DeepCopy()
		if secret.Labels == nil {
			secret.Labels = map[string]string{}
		}
		secret.Labels[planner.MachineNamespaceLabel] = capiCluster.Namespace
		secret.Labels[planner.MachineNameLabel] = secret.Name

		return h.secrets.Update(secret)
	}

	return secret, nil
}

func (h *handler) createMachine(capiCluster *capi.Cluster, secret *corev1.Secret, data data.Object) error {
	objs, err := h.createMachineObjects(capiCluster, secret.Name, data)
	if err != nil {
		return err
	}
	return h.apply.WithOwner(secret).ApplyObjects(objs...)
}

func (h *handler) createMachineObjects(capiCluster *capi.Cluster, machineName string, data data.Object) ([]runtime.Object, error) {
	labels := map[string]string{}
	annotations := map[string]string{}

	if data.Bool("role-control-plane") {
		labels[planner.ControlPlaneRoleLabel] = "true"
	}
	if data.Bool("role-etcd") {
		labels[planner.EtcdRoleLabel] = "true"
	}
	if data.Bool("role-worker") {
		labels[planner.WorkerRoleLabel] = "true"
	}

	labelsMap := map[string]string{}
	for _, str := range strings.Split(data.String("label"), ",") {
		k, v := kv.Split(str, "=")
		if k == "" {
			continue
		}
		labelsMap[k] = v
	}

	if len(labelsMap) > 0 {
		data, err := json.Marshal(labelsMap)
		if err != nil {
			return nil, err
		}
		annotations[planner.LabelsAnnotation] = string(data)
	}

	var taints []corev1.Taint
	for _, taint := range convert.ToStringSlice(data["taints"]) {
		for _, taint := range strings.Split(taint, ",") {
			parts := strings.Split(taint, ":")
			switch len(parts) {
			case 1:
				taints = append(taints, corev1.Taint{
					Key: parts[0],
				})
			case 2:
				taints = append(taints, corev1.Taint{
					Key:   parts[0],
					Value: parts[1],
				})
			case 3:
				taints = append(taints, corev1.Taint{
					Key:    parts[0],
					Value:  parts[1],
					Effect: corev1.TaintEffect(parts[2]),
				})
			}
		}
	}

	if len(taints) > 0 {
		data, err := json.Marshal(taints)
		if err != nil {
			return nil, err
		}
		annotations[planner.TaintsAnnotation] = string(data)
	}

	return []runtime.Object{
		&rkev1.RKEBootstrap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      machineName,
				Namespace: capiCluster.Namespace,
			},
		},
		&rkev1.UnmanagedMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      machineName,
				Namespace: capiCluster.Namespace,
			},
		},
		&capi.Machine{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:        machineName,
				Namespace:   capiCluster.Namespace,
				Labels:      labels,
				Annotations: annotations,
			},
			Spec: capi.MachineSpec{
				ClusterName: capiCluster.Name,
				Bootstrap: capi.Bootstrap{
					ConfigRef: &corev1.ObjectReference{
						Kind:       "RKEBootstrap",
						Namespace:  capiCluster.Namespace,
						Name:       machineName,
						APIVersion: "rke.cattle.io/v1",
					},
				},
				InfrastructureRef: corev1.ObjectReference{
					Kind:       "UnmanagedMachine",
					Namespace:  capiCluster.Namespace,
					Name:       machineName,
					APIVersion: "rke.cattle.io/v1",
				},
			},
		},
	}, nil
}

func (h *handler) getCAPICluster(secret *corev1.Secret) (*capi.Cluster, error) {
	cluster, err := h.mgmtClusterCache.Get(secret.Namespace)
	if apierror.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	owner, err := h.apply.FindOwner(cluster)
	if err == apply.ErrOwnerNotFound {
		return nil, nil
	}

	rcluster, ok := owner.(*rancherv1.Cluster)
	if !ok {
		return nil, nil
	}

	return h.capiClusterCache.Get(rcluster.Namespace, rcluster.Name)
}

func (h *handler) onUnmanagedMachineChange(key string, machine *rkev1.UnmanagedMachine) (*rkev1.UnmanagedMachine, error) {
	if machine != nil && !machine.Status.Ready {
		machine = machine.DeepCopy()
		machine.Status.Ready = true
		return h.unmanagedMachine.UpdateStatus(machine)
	}
	return machine, nil
}
