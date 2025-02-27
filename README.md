# IstiodPOCBluepi
I demo service to learn more about istiod
# Microservices Project with Istio and Datadog Monitoring
## https://app.diagrams.net/#G1e8vDn2fPGrLEDZYzFHExLPR-PIuoWTvC#%7B%22pageId%22%3A%225CSxC_CJXj1vpsv0aLtl%22%7D
## Overview
This project consists of three Golang-based microservices deployed in a Kubernetes cluster with Istio for service-to-service security and observability. The setup includes JWT authentication via Auth0 and centralized logging using Persistent Volumes and Datadog.

## Architecture
```
┌──────────────┐        ┌──────────────┐        ┌──────────────┐
│   Service A  │  --->  │   Service C  │  --->  │   Service B  │
└──────────────┘        └──────────────┘        └──────────────┘
        │                      │                      │
        ▼                      ▼                      ▼
    Istio Ingress          Istio Sidecar           MySQL Database
    Gateway (mTLS)            (mTLS)
```

- **Service A**: Handles user authentication (Auth0) and routes requests to Service C.
- **Service C**: Acts as an intermediary, validates JWT tokens, and forwards requests to Service B.
- **Service B**: Connects to MySQL and processes requests.
- **Istio**: Enables mTLS, logging, and monitoring.
- **Datadog**: Monitors Kubernetes metrics, logs, and traces.

## Deployment
### Prerequisites
Ensure you have the following installed:
- **Kubernetes** (Minikube or a cloud-based cluster)
- **Kubectl**
- **Istio** (installed via Helm)
- **Helm** (for managing charts)
- **Datadog Agent** (for monitoring)

### Deploy Services
Apply Kubernetes manifests to deploy the services:
```sh
kubectl apply -f k8s/service-a.yaml
kubectl apply -f k8s/service-b.yaml
kubectl apply -f k8s/service-c.yaml
```

### Expose via Istio Ingress Gateway
```sh
kubectl apply -f k8s/istio-gateway.yaml
```

### Apply Istio Telemetry for Logging
```sh
kubectl apply -f k8s/istio-telemetry.yaml
```

### Install Datadog via Helm
```sh
helm install datadog-agent datadog/datadog --namespace default --values dd-values.yaml
```

## Logging Setup
### Persistent Volume for Logs
Each service logs to a Persistent Volume to retain logs even after pod restarts.

#### Example of Service A Logging Setup:
```yaml
volumes:
  - name: istio-storage
    persistentVolumeClaim:
      claimName: istio-a-pvc
volumeMounts:
  - name: istio-storage
    mountPath: /var/log/istio
```

### Viewing Logs
#### View Service Logs
```sh
kubectl exec -it <pod-name> -c service-a -- tail -f /var/log/service-a/service-a.log
```

#### View Istio Proxy Logs
```sh
kubectl exec -it <pod-name> -c istio-proxy -- tail -f /var/log/istio/sidecar-a-access.log
```

#### View Datadog Agent Logs
```sh
kubectl logs -l app=datadog-agent
```

## Monitoring with Datadog
Datadog monitors logs, traces, and metrics from the Kubernetes cluster and Istio mesh.
- View metrics on the Datadog dashboard.
- Use `kubectl logs` to check Datadog agent logs for troubleshooting.

## Troubleshooting
### Pods Not Running?
```sh
kubectl get pods -n default
```
Check logs for errors:
```sh
kubectl logs <pod-name>
```

### Helm Upgrade Fails?
Ensure the correct release name:
```sh
helm upgrade datadog-agent datadog/datadog --namespace default --values dd-values.yaml
```

## Conclusion
This project demonstrates a secure microservices architecture using Istio and Auth0 for authentication, with Datadog for monitoring and centralized logging. 🚀

