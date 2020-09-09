# NetworkPolicy Operator

This operator extends [networking.k8s.io/v1/NetworkPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#networkpolicy-v1-networking-k8s-io) functionnalities by adding a field `.spec.egress[].to[].domain`. An underlying `networking.k8s.io/v1/NetworkPolicy` is created with the IP addresses resolved from specified domains every `.spec.resolveEverySeconds` (default to 300 seconds).

Note that this is not ready for production and this practice have a big downside:

DNS servers may respond differently based on your location, your IP, etc... thus the pods may resolve domain names to different addresses and have their traffic denied, which is somewhat very inconvenient.

An alternative is to use a service mesh like Istio (cf. https://istio.io/latest/docs/tasks/traffic-management/egress/).
