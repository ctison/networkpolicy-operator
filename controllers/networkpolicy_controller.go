package controllers

import (
	"context"
	"net"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/k7o-io/networkpolicy-operator/api/v1alpha1"
)

// NetworkPolicyReconciler reconciles a NetworkPolicy object
type NetworkPolicyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=networking.k7o.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k7o.io,resources=networkpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies/status,verbs=get;list;watch;create;update;patch;delete

func (r *NetworkPolicyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("networkpolicy", req.NamespacedName)

	// Fetch the NetworkPolicy resource
	networkPolicy := &v1alpha1.NetworkPolicy{}
	if err := r.Get(ctx, req.NamespacedName, networkPolicy); err != nil {
		if errors.IsNotFound(err) {
			log.Info("NetworkPolicy resource not found... skipping")
			return ctrl.Result{}, nil
		}
		log.Error(err, "failed to fetch NetworkPolicy resource")
		return ctrl.Result{}, err
	}

	// Compute underlying NetworkPolicy
	state := &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "NetworkPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      networkPolicy.Name,
			Namespace: networkPolicy.Namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			Egress:      make([]networkingv1.NetworkPolicyEgressRule, len(networkPolicy.Spec.Egress)),
			Ingress:     networkPolicy.Spec.Ingress,
			PodSelector: networkPolicy.Spec.PodSelector,
			PolicyTypes: networkPolicy.Spec.PolicyTypes,
		},
	}
	for i, egress := range networkPolicy.Spec.Egress {
		state.Spec.Egress[i].Ports = egress.Ports
		if egress.To == nil {
			continue
		}
		state.Spec.Egress[i].To = make([]networkingv1.NetworkPolicyPeer, 0, len(egress.To))
		for _, to := range egress.To {
			if to.Domain == nil {
				state.Spec.Egress[i].To = append(state.Spec.Egress[i].To, networkingv1.NetworkPolicyPeer{
					NamespaceSelector: to.NamespaceSelector,
					IPBlock:           to.IPBlock,
					PodSelector:       to.PodSelector,
				})
				continue
			}
			ipAddrs, err := net.LookupIP(*to.Domain)
			if err != nil {
				log.Error(err, "cannot resolve domain... skipping", "domain", *to.Domain)
				continue
			}
			if len(ipAddrs) == 0 {
				log.Info("domain resolves to no IP... skipping", "domain", *to.Domain)
				continue
			}
			for _, ipAddr := range ipAddrs {
				state.Spec.Egress[i].To = append(state.Spec.Egress[i].To, networkingv1.NetworkPolicyPeer{
					NamespaceSelector: to.NamespaceSelector,
					IPBlock: &networkingv1.IPBlock{
						CIDR: ipAddr.String() + "/32",
					},
					PodSelector: to.PodSelector,
				})
			}
		}
	}
	if err := ctrl.SetControllerReference(networkPolicy, state, r.Scheme); err != nil {
		log.Error(err, "!!!PANIC cannot set controller reference", "namespace", state.Namespace, "name", state.Name)
		return ctrl.Result{}, err
	}

	// Check if NetworkPolicy already exists, if not create a new one
	found := &networkingv1.NetworkPolicy{}
	err := r.Get(ctx, types.NamespacedName{Name: state.Name, Namespace: state.Namespace}, found)
	if err != nil {
		if !errors.IsNotFound(err) {
			log.Error(err, "failed to get NetworkPolicy", "namespace", state.Namespace, "name", state.Name)
			return ctrl.Result{}, err
		}
		log.Info("creating a new NetworkPolicy", "namespace", state.Namespace, "name", state.Name)
		if err = r.Create(ctx, state); err != nil {
			log.Error(err, "failed to create a new NetworkPolicy", "namespace", state.Namespace, "name", state.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Ensure underlying NetworkPolicy is up to date
	if !reflect.DeepEqual(found, state) {
		if err := r.Update(ctx, state); err != nil {
			log.Error(err, "failed to update NetworkPolicy", "namespace", found.Namespace, "name", found.Name)
			return ctrl.Result{}, err
		}
		log.Info("successfully updated NetworkPolicy", "namespace", found.Namespace, "name", found.Name)
		return ctrl.Result{Requeue: true}, nil
	}

	// Update status. if needed
	if networkPolicy.Status.NetworkPolicyName == nil || *networkPolicy.Status.NetworkPolicyName != state.Name {
		networkPolicy.Status.NetworkPolicyName = &state.Name
		if err := r.Status().Update(ctx, networkPolicy); err != nil {
			log.Error(err, "Failed to update NetworkPolicy status", "namespace", networkPolicy.Namespace, "name", networkPolicy.Name)
			return ctrl.Result{}, err
		}
	}

	resolveEverySeconds := v1alpha1.DefaultResolveEverySeconds
	if networkPolicy.Spec.ResolveEverySeconds != nil {
		resolveEverySeconds = *networkPolicy.Spec.ResolveEverySeconds
	}
	return ctrl.Result{RequeueAfter: time.Duration(resolveEverySeconds) * time.Second}, nil
}

func (r *NetworkPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NetworkPolicy{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Complete(r)
}

func (r *NetworkPolicyReconciler) makeNetworkPolicy(from *v1alpha1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{
		ObjectMeta: from.ObjectMeta,
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: from.Spec.PodSelector,
			Ingress:     from.Spec.Ingress,
			PolicyTypes: from.Spec.PolicyTypes,
			Egress:      nil,
		},
	}
	// egress := make([]networkingv1.NetworkPolicyEgressRule, 0, len(from.Spec.Egress))
	return netpol, nil
}
