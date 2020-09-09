package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var networkpolicylog = logf.Log.WithName("networkpolicy-resource")

func (r *NetworkPolicy) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-networking-k7o-io-v1alpha1-networkpolicy,mutating=true,failurePolicy=fail,groups=networking.k7o.io,resources=networkpolicies,verbs=create;update,versions=v1alpha1,name=mnetworkpolicy.kb.io

var _ webhook.Defaulter = &NetworkPolicy{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *NetworkPolicy) Default() {
	networkpolicylog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-networking-k7o-io-v1alpha1-networkpolicy,mutating=false,failurePolicy=fail,groups=networking.k7o.io,resources=networkpolicies,versions=v1alpha1,name=vnetworkpolicy.kb.io

var _ webhook.Validator = &NetworkPolicy{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NetworkPolicy) ValidateCreate() error {
	networkpolicylog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NetworkPolicy) ValidateUpdate(old runtime.Object) error {
	networkpolicylog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NetworkPolicy) ValidateDelete() error {
	networkpolicylog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
