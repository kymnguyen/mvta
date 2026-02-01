# Deployment Pipeline Proposal: GitOps, ArgoCD, and Kubernetes

**Project:** MVTA (Multi-Vehicle Tracking Application)  
**Date:** February 1, 2026  
**Version:** 1.0

---

## Executive Summary

This proposal outlines a modern, automated deployment pipeline for the MVTA platform using GitOps principles, ArgoCD as the continuous delivery tool, and Kubernetes as the orchestration platform. This approach will enable:

- **Automated deployments** with declarative infrastructure
- **Enhanced reliability** through immutable infrastructure patterns
- **Improved security** with audit trails and RBAC
- **Faster time-to-market** with streamlined CI/CD workflows
- **Easy rollbacks** and disaster recovery

---

## 1. Current State Analysis

### Existing Architecture
- **Backend Services:** 4 Go microservices (auth-svc, tracking-svc, vehicle-svc, workflow-svc)
- **Frontend:** React-based admin web application
- **Infrastructure:** Docker Compose for local development
- **Deployment:** Manual or script-based deployments

### Challenges
- No standardized deployment process across environments
- Limited deployment automation
- Manual configuration management
- Difficulty tracking configuration changes
- No self-healing infrastructure
- Complex rollback procedures

---

## 2. Proposed Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Git Repository                          │
│  ┌──────────────────┐         ┌──────────────────┐         │
│  │  Application     │         │   Infrastructure  │         │
│  │  Source Code     │         │   as Code (K8s)  │         │
│  └──────────────────┘         └──────────────────┘         │
└────────┬──────────────────────────────┬───────────────────┘
         │                               │
         │ Push                          │ Sync
         ▼                               ▼
┌─────────────────────┐        ┌──────────────────────┐
│   CI Pipeline       │        │      ArgoCD          │
│   (GitHub Actions)  │        │   (GitOps Engine)    │
│                     │        │                      │
│ • Build             │        │ • Monitors Git       │
│ • Test              │        │ • Applies Changes    │
│ • Scan              │────────▶• Self-Heals          │
│ • Push Images       │        │ • Syncs State        │
└─────────────────────┘        └──────────┬───────────┘
                                          │ Deploy
                                          ▼
                              ┌──────────────────────┐
                              │  Kubernetes Cluster  │
                              │                      │
                              │  ┌────────────────┐ │
                              │  │   Namespaces   │ │
                              │  │ • dev          │ │
                              │  │ • staging      │ │
                              │  │ • production   │ │
                              │  └────────────────┘ │
                              └──────────────────────┘
```

### Component Breakdown

#### 2.1 Git Repository Structure
```
mvta/
├── apps/                           # Application source code
│   ├── admin-web/
│   └── backend/
├── infra/
│   ├── k8s/
│   │   ├── base/                   # Base Kubernetes manifests
│   │   │   ├── auth-svc/
│   │   │   │   ├── deployment.yaml
│   │   │   │   ├── service.yaml
│   │   │   │   └── configmap.yaml
│   │   │   ├── tracking-svc/
│   │   │   ├── vehicle-svc/
│   │   │   ├── workflow-svc/
│   │   │   ├── admin-web/
│   │   │   └── shared/             # Shared resources (ingress, etc.)
│   │   └── overlays/               # Environment-specific configs
│   │       ├── dev/
│   │       │   └── kustomization.yaml
│   │       ├── staging/
│   │       │   └── kustomization.yaml
│   │       └── production/
│   │           └── kustomization.yaml
│   └── argocd/
│       ├── applications/            # ArgoCD Application definitions
│       │   ├── mvta-dev.yaml
│       │   ├── mvta-staging.yaml
│       │   └── mvta-production.yaml
│       └── projects/
│           └── mvta-project.yaml
└── .github/
    └── workflows/
        ├── ci-backend.yaml
        ├── ci-frontend.yaml
        └── image-promote.yaml
```

#### 2.2 Kubernetes Cluster Architecture
- **Multi-namespace setup** for environment isolation
- **RBAC policies** for security
- **Network policies** for service communication
- **Resource quotas** and limits per namespace
- **Horizontal Pod Autoscaling** for dynamic scaling

#### 2.3 ArgoCD Configuration
- **Application-per-microservice** pattern
- **App of Apps** pattern for managing multiple applications
- **Automated sync policies** with self-healing
- **Progressive sync strategies** with health checks
- **Multi-environment management** (dev, staging, production)

---

## 3. GitOps Principles

### Core Tenets

1. **Declarative Configuration**
   - All infrastructure and application configs defined declaratively in Git
   - Kubernetes manifests using Kustomize for environment-specific overlays

2. **Version Control as Source of Truth**
   - Git repository is the single source of truth
   - All changes tracked with full audit history
   - Easy rollback to any previous state

3. **Automated Deployment**
   - ArgoCD continuously monitors Git repository
   - Automatic synchronization of desired vs. actual state
   - Self-healing when drift is detected

4. **Observable and Verifiable**
   - Real-time visibility into deployment status
   - Automated health checks and validation
   - Drift detection and alerts

---

## 4. CI/CD Pipeline Design

### 4.1 Continuous Integration (CI)

**Trigger:** Push to feature branch or pull request

**Pipeline Stages:**

```yaml
# .github/workflows/ci-backend.yaml
name: Backend CI Pipeline

on:
  push:
    paths:
      - 'apps/backend/**'
  pull_request:
    paths:
      - 'apps/backend/**'

jobs:
  build-and-test:
    strategy:
      matrix:
        service: [auth-svc, tracking-svc, vehicle-svc, workflow-svc]
    
    steps:
      # 1. Checkout code
      # 2. Setup Go environment
      # 3. Run tests and linting
      # 4. Build binary
      # 5. Security scanning (Trivy, Snyk)
      # 6. Build Docker image
      # 7. Push to container registry (with commit SHA tag)
      # 8. Update image tag in Git (for dev environment)
```

**Stages:**
1. **Lint & Format Check** - Go linting, formatting validation
2. **Unit Tests** - Run all unit tests with coverage reports
3. **Integration Tests** - Test service interactions
4. **Security Scanning** - Scan for vulnerabilities (Trivy, Snyk)
5. **Build Container Image** - Multi-stage Docker builds
6. **Image Scanning** - Scan container images for CVEs
7. **Push to Registry** - Tag with commit SHA and push to registry
8. **Update Manifest** - Update image tags in Git repository

### 4.2 Continuous Deployment (CD)

**Trigger:** Image tag update in k8s manifests

**ArgoCD Workflow:**

```yaml
# infra/argocd/applications/mvta-dev.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mvta-dev
  namespace: argocd
spec:
  project: mvta
  source:
    repoURL: https://github.com/kymnguyen/mvta.git
    targetRevision: main
    path: infra/k8s/overlays/dev
  destination:
    server: https://kubernetes.default.svc
    namespace: mvta-dev
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
```

**ArgoCD Deployment Flow:**
1. **Detect Change** - ArgoCD polls Git repository every 3 minutes
2. **Compare State** - Compare desired (Git) vs actual (K8s) state
3. **Sync Resources** - Apply changes to Kubernetes cluster
4. **Health Check** - Verify pod health and readiness
5. **Rollout Strategy** - Rolling update with zero downtime
6. **Self-Heal** - Automatically correct any drift

### 4.3 Promotion Strategy

**Image Promotion Pipeline:**

```
Dev Environment (auto-deploy)
    ↓ (manual promotion after validation)
Staging Environment (manual approval)
    ↓ (manual promotion after testing)
Production Environment (manual approval + deployment window)
```

**Implementation:**
```yaml
# .github/workflows/image-promote.yaml
name: Promote Image to Environment

on:
  workflow_dispatch:
    inputs:
      service:
        description: 'Service to promote'
        required: true
        type: choice
        options:
          - auth-svc
          - tracking-svc
          - vehicle-svc
          - workflow-svc
      image_tag:
        description: 'Image tag to promote'
        required: true
      target_environment:
        description: 'Target environment'
        required: true
        type: choice
        options:
          - staging
          - production

jobs:
  promote:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      
      - name: Update image tag
        run: |
          cd infra/k8s/overlays/${{ inputs.target_environment }}
          kustomize edit set image ${{ inputs.service }}=registry/mvta/${{ inputs.service }}:${{ inputs.image_tag }}
      
      - name: Commit and push
        run: |
          git config user.name "GitHub Actions"
          git config user.email "actions@github.com"
          git add .
          git commit -m "Promote ${{ inputs.service }} to ${{ inputs.target_environment }}: ${{ inputs.image_tag }}"
          git push
```

---

## 5. Kubernetes Configuration

### 5.1 Base Manifests Example

**Deployment Manifest:**
```yaml
# infra/k8s/base/auth-svc/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auth-svc
  labels:
    app: auth-svc
    version: v1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: auth-svc
  template:
    metadata:
      labels:
        app: auth-svc
        version: v1
    spec:
      containers:
      - name: auth-svc
        image: registry.example.com/mvta/auth-svc:latest
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 9090
          name: metrics
        env:
        - name: ENV
          value: "dev"
        - name: PORT
          value: "8080"
        envFrom:
        - configMapRef:
            name: auth-svc-config
        - secretRef:
            name: auth-svc-secrets
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
```

**Service Manifest:**
```yaml
# infra/k8s/base/auth-svc/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: auth-svc
  labels:
    app: auth-svc
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    protocol: TCP
    name: http
  - port: 9090
    targetPort: 9090
    protocol: TCP
    name: metrics
  selector:
    app: auth-svc
```

### 5.2 Kustomize Overlays

**Development Environment:**
```yaml
# infra/k8s/overlays/dev/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: mvta-dev

bases:
  - ../../base/auth-svc
  - ../../base/tracking-svc
  - ../../base/vehicle-svc
  - ../../base/workflow-svc
  - ../../base/admin-web
  - ../../base/shared

replicas:
  - name: auth-svc
    count: 1
  - name: tracking-svc
    count: 1
  - name: vehicle-svc
    count: 1
  - name: workflow-svc
    count: 1

images:
  - name: registry.example.com/mvta/auth-svc
    newTag: dev-latest
  - name: registry.example.com/mvta/tracking-svc
    newTag: dev-latest
  - name: registry.example.com/mvta/vehicle-svc
    newTag: dev-latest
  - name: registry.example.com/mvta/workflow-svc
    newTag: dev-latest

configMapGenerator:
  - name: global-config
    literals:
      - ENV=dev
      - LOG_LEVEL=debug
      - ENABLE_METRICS=true

secretGenerator:
  - name: global-secrets
    literals:
      - JWT_SECRET=dev-secret-key
    type: Opaque
```

**Production Environment:**
```yaml
# infra/k8s/overlays/production/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: mvta-production

bases:
  - ../../base/auth-svc
  - ../../base/tracking-svc
  - ../../base/vehicle-svc
  - ../../base/workflow-svc
  - ../../base/admin-web
  - ../../base/shared

replicas:
  - name: auth-svc
    count: 3
  - name: tracking-svc
    count: 3
  - name: vehicle-svc
    count: 3
  - name: workflow-svc
    count: 2

images:
  - name: registry.example.com/mvta/auth-svc
    newTag: v1.0.0
  - name: registry.example.com/mvta/tracking-svc
    newTag: v1.0.0
  - name: registry.example.com/mvta/vehicle-svc
    newTag: v1.0.0
  - name: registry.example.com/mvta/workflow-svc
    newTag: v1.0.0

patchesStrategicMerge:
  - production-resources.yaml
  - production-hpa.yaml

configMapGenerator:
  - name: global-config
    literals:
      - ENV=production
      - LOG_LEVEL=info
      - ENABLE_METRICS=true

secretGenerator:
  - name: global-secrets
    files:
      - secrets/jwt-secret.txt
    type: Opaque
```

### 5.3 Advanced Kubernetes Resources

**Horizontal Pod Autoscaler:**
```yaml
# infra/k8s/overlays/production/production-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: auth-svc-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: auth-svc
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

**Network Policy:**
```yaml
# infra/k8s/base/shared/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: auth-svc-network-policy
spec:
  podSelector:
    matchLabels:
      app: auth-svc
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: admin-web
    - podSelector:
        matchLabels:
          app: tracking-svc
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - protocol: TCP
      port: 5432
```

---

## 6. ArgoCD Implementation

### 6.1 Installation

**Namespace and Installation:**
```bash
# Create namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Expose ArgoCD server (for testing)
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

### 6.2 Project Configuration

```yaml
# infra/argocd/projects/mvta-project.yaml
apiVersion: argoproj.io/v1alpha1
kind: AppProject
metadata:
  name: mvta
  namespace: argocd
spec:
  description: MVTA Vehicle Tracking Platform
  
  sourceRepos:
    - https://github.com/kymnguyen/mvta.git
  
  destinations:
    - namespace: mvta-dev
      server: https://kubernetes.default.svc
    - namespace: mvta-staging
      server: https://kubernetes.default.svc
    - namespace: mvta-production
      server: https://kubernetes.default.svc
  
  clusterResourceWhitelist:
    - group: ''
      kind: Namespace
  
  namespaceResourceWhitelist:
    - group: 'apps'
      kind: Deployment
    - group: ''
      kind: Service
    - group: ''
      kind: ConfigMap
    - group: ''
      kind: Secret
    - group: 'networking.k8s.io'
      kind: Ingress
    - group: 'autoscaling'
      kind: HorizontalPodAutoscaler
  
  roles:
    - name: developer
      description: Developers can sync dev environment
      policies:
        - p, proj:mvta:developer, applications, sync, mvta/mvta-dev, allow
        - p, proj:mvta:developer, applications, get, mvta/*, allow
    
    - name: operator
      description: Operators can sync staging and production
      policies:
        - p, proj:mvta:operator, applications, *, mvta/*, allow
```

### 6.3 Application Definitions

**App of Apps Pattern:**
```yaml
# infra/argocd/applications/root-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mvta-root
  namespace: argocd
spec:
  project: mvta
  source:
    repoURL: https://github.com/kymnguyen/mvta.git
    targetRevision: main
    path: infra/argocd/applications
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

**Individual Application:**
```yaml
# infra/argocd/applications/mvta-dev.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mvta-dev
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: mvta
  
  source:
    repoURL: https://github.com/kymnguyen/mvta.git
    targetRevision: main
    path: infra/k8s/overlays/dev
  
  destination:
    server: https://kubernetes.default.svc
    namespace: mvta-dev
  
  syncPolicy:
    automated:
      prune: true          # Remove resources not in Git
      selfHeal: true       # Auto-sync if cluster state drifts
      allowEmpty: false    # Prevent deletion of all resources
    
    syncOptions:
      - CreateNamespace=true
      - PrunePropagationPolicy=foreground
      - PruneLast=true
    
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m
  
  ignoreDifferences:
    - group: apps
      kind: Deployment
      jsonPointers:
        - /spec/replicas  # Ignore replica count (managed by HPA)
```

**Production with Manual Sync:**
```yaml
# infra/argocd/applications/mvta-production.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mvta-production
  namespace: argocd
  annotations:
    notifications.argoproj.io/subscribe.on-sync-succeeded.slack: mvta-alerts
spec:
  project: mvta
  
  source:
    repoURL: https://github.com/kymnguyen/mvta.git
    targetRevision: main
    path: infra/k8s/overlays/production
  
  destination:
    server: https://kubernetes.default.svc
    namespace: mvta-production
  
  syncPolicy:
    # No automated sync for production - manual approval required
    syncOptions:
      - CreateNamespace=true
    
    retry:
      limit: 3
      backoff:
        duration: 10s
        factor: 2
        maxDuration: 5m
  
  # Health assessment
  health:
    - group: apps
      kind: Deployment
      check: |
        hs = {}
        if obj.status ~= nil then
          if obj.status.updatedReplicas == obj.spec.replicas then
            hs.status = "Healthy"
            hs.message = "All replicas are updated"
            return hs
          end
        end
        hs.status = "Progressing"
        hs.message = "Waiting for rollout to finish"
        return hs
```

### 6.4 Progressive Delivery with Argo Rollouts

**Canary Deployment Strategy:**
```yaml
# infra/k8s/base/auth-svc/rollout.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: auth-svc
spec:
  replicas: 5
  strategy:
    canary:
      steps:
      - setWeight: 20
      - pause: {duration: 5m}
      - setWeight: 40
      - pause: {duration: 5m}
      - setWeight: 60
      - pause: {duration: 5m}
      - setWeight: 80
      - pause: {duration: 5m}
      
      canaryService: auth-svc-canary
      stableService: auth-svc-stable
      
      trafficRouting:
        istio:
          virtualService:
            name: auth-svc
            routes:
            - primary
      
      analysis:
        templates:
        - templateName: success-rate
        startingStep: 2
        args:
        - name: service-name
          value: auth-svc
  
  selector:
    matchLabels:
      app: auth-svc
  
  template:
    metadata:
      labels:
        app: auth-svc
    spec:
      containers:
      - name: auth-svc
        image: registry.example.com/mvta/auth-svc:latest
        # ... rest of container spec
```

---

## 7. Security Considerations

### 7.1 Secrets Management

**Sealed Secrets:**
```bash
# Install Sealed Secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Encrypt a secret
kubectl create secret generic auth-db-credentials \
  --from-literal=username=admin \
  --from-literal=password=secretpass \
  --dry-run=client -o yaml | \
  kubeseal -o yaml > sealed-secret.yaml
```

**External Secrets Operator (Alternative):**
```yaml
# Using AWS Secrets Manager / Azure Key Vault / GCP Secret Manager
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: auth-svc-secrets
spec:
  refreshInterval: 1h
  secretStoreRef:
    name: aws-secrets-manager
    kind: SecretStore
  target:
    name: auth-svc-secrets
    creationPolicy: Owner
  data:
    - secretKey: database_url
      remoteRef:
        key: mvta/production/auth-svc/database-url
    - secretKey: jwt_secret
      remoteRef:
        key: mvta/production/auth-svc/jwt-secret
```

### 7.2 RBAC Configuration

**Service Account for ArgoCD:**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argocd-application-controller
  namespace: argocd
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: argocd-application-controller
rules:
- apiGroups: ['*']
  resources: ['*']
  verbs: ['get', 'list', 'watch']
- apiGroups: ['apps']
  resources: ['deployments', 'replicasets']
  verbs: ['get', 'list', 'watch', 'update', 'patch', 'delete']
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argocd-application-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: argocd-application-controller
subjects:
- kind: ServiceAccount
  name: argocd-application-controller
  namespace: argocd
```

### 7.3 Security Scanning in CI

```yaml
# .github/workflows/security-scan.yaml
name: Security Scanning

jobs:
  trivy-scan:
    runs-on: ubuntu-latest
    steps:
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: registry.example.com/mvta/auth-svc:${{ github.sha }}
          format: 'sarif'
          output: 'trivy-results.sarif'
          severity: 'CRITICAL,HIGH'
      
      - name: Upload to Security tab
        uses: github/codeql-action/upload-sarif@v2
        with:
          sarif_file: 'trivy-results.sarif'
  
  snyk-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Snyk to check for vulnerabilities
        uses: snyk/actions/golang@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
        with:
          command: test
          args: --severity-threshold=high
```

### 7.4 Network Security

**Service Mesh Integration (Istio):**
```yaml
# Enable mTLS between services
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: mvta-production
spec:
  mtls:
    mode: STRICT
```

---

## 8. Monitoring and Observability

### 8.1 Prometheus + Grafana Stack

**Prometheus ServiceMonitor:**
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: auth-svc-metrics
  namespace: mvta-production
spec:
  selector:
    matchLabels:
      app: auth-svc
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

**Grafana Dashboard ConfigMap:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mvta-dashboard
  namespace: monitoring
data:
  mvta-overview.json: |
    {
      "dashboard": {
        "title": "MVTA Platform Overview",
        "panels": [
          {
            "title": "Request Rate",
            "targets": [
              {
                "expr": "sum(rate(http_requests_total[5m])) by (service)"
              }
            ]
          }
        ]
      }
    }
```

### 8.2 Logging with EFK Stack

**Fluentd DaemonSet:**
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: logging
spec:
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      containers:
      - name: fluentd
        image: fluent/fluentd-kubernetes-daemonset:v1-debian-elasticsearch
        env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch.logging.svc.cluster.local"
        - name: FLUENT_ELASTICSEARCH_PORT
          value: "9200"
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: varlibdockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: varlibdockercontainers
        hostPath:
          path: /var/lib/docker/containers
```

### 8.3 ArgoCD Notifications

```yaml
# infra/argocd/notifications-cm.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    token: $slack-token
  
  template.app-deployed: |
    message: |
      Application {{.app.metadata.name}} is now running new version.
      Service: {{.app.spec.source.path}}
      Revision: {{.app.status.sync.revision}}
  
  template.app-health-degraded: |
    message: |
      Application {{.app.metadata.name}} has degraded health status.
      Details: {{.app.status.health.message}}
  
  trigger.on-deployed: |
    - when: app.status.operationState.phase in ['Succeeded']
      send: [app-deployed]
  
  trigger.on-health-degraded: |
    - when: app.status.health.status == 'Degraded'
      send: [app-health-degraded]
```

---

## 9. Disaster Recovery and Rollback

### 9.1 Rollback Strategy

**Manual Rollback via Git:**
```bash
# Revert to previous commit
git revert HEAD
git push origin main

# ArgoCD will automatically sync to the previous state
```

**ArgoCD CLI Rollback:**
```bash
# View application history
argocd app history mvta-production

# Rollback to specific revision
argocd app rollback mvta-production 10
```

### 9.2 Backup Strategy

**Velero for Cluster Backups:**
```bash
# Install Velero
velero install \
  --provider aws \
  --plugins velero/velero-plugin-for-aws:v1.8.0 \
  --bucket mvta-k8s-backups \
  --backup-location-config region=us-east-1

# Create backup schedule
velero schedule create mvta-daily \
  --schedule="0 2 * * *" \
  --include-namespaces mvta-production

# Restore from backup
velero restore create --from-backup mvta-daily-20260201
```

---

## 10. Migration Strategy

### Phase 1: Infrastructure Setup (Week 1-2)
- [ ] Provision Kubernetes cluster (EKS/AKS/GKE)
- [ ] Install ArgoCD in cluster
- [ ] Set up container registry
- [ ] Configure networking and ingress

### Phase 2: CI Pipeline Implementation (Week 2-3)
- [ ] Create GitHub Actions workflows for all services
- [ ] Implement automated testing
- [ ] Set up security scanning
- [ ] Configure image building and pushing

### Phase 3: Kubernetes Manifest Creation (Week 3-4)
- [ ] Create base Kubernetes manifests for all services
- [ ] Set up Kustomize overlays for dev/staging/prod
- [ ] Configure ConfigMaps and Secrets
- [ ] Define resource limits and HPA policies

### Phase 4: ArgoCD Configuration (Week 4-5)
- [ ] Create ArgoCD projects and applications
- [ ] Configure sync policies
- [ ] Set up RBAC and permissions
- [ ] Enable notifications

### Phase 5: Dev Environment Deployment (Week 5)
- [ ] Deploy to dev namespace
- [ ] Validate all services are running
- [ ] Test inter-service communication
- [ ] Verify automated sync works

### Phase 6: Staging Environment (Week 6)
- [ ] Deploy to staging namespace
- [ ] Run end-to-end tests
- [ ] Performance testing
- [ ] Security testing

### Phase 7: Production Readiness (Week 7)
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy
- [ ] Document runbooks
- [ ] Train team on new processes

### Phase 8: Production Migration (Week 8)
- [ ] Blue-green deployment to production
- [ ] Monitor closely for issues
- [ ] Validate all functionality
- [ ] Update DNS/routing

### Phase 9: Post-Migration (Week 9+)
- [ ] Decommission old infrastructure
- [ ] Optimize resource usage
- [ ] Implement progressive delivery
- [ ] Continuous improvement

---

## 11. Benefits

### Technical Benefits
- **Automated deployments** reducing manual errors
- **Self-healing infrastructure** with automatic drift correction
- **Declarative configuration** making infrastructure auditable
- **Easy rollbacks** to any previous state
- **Progressive delivery** with canary and blue-green deployments
- **Multi-environment consistency** with Kustomize overlays

### Operational Benefits
- **Faster time-to-market** with automated pipelines
- **Improved reliability** through automated testing and validation
- **Better visibility** with centralized GitOps dashboard
- **Reduced cognitive load** on operators
- **Audit trail** of all infrastructure changes

### Business Benefits
- **Lower operational costs** through automation
- **Increased development velocity** 
- **Better compliance** with change tracking
- **Reduced downtime** with sophisticated deployment strategies
- **Scalability** to handle growth

---

## 12. Risks and Mitigation

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| **Learning curve for team** | Medium | High | Comprehensive training, documentation, and phased rollout |
| **Initial setup complexity** | Medium | Medium | Use experienced consultants, follow best practices |
| **ArgoCD downtime** | High | Low | Run ArgoCD in HA mode, maintain runbooks for manual deployment |
| **Kubernetes cluster failure** | High | Low | Multi-zone deployment, regular backups with Velero |
| **Git repository compromise** | High | Low | Branch protection, signed commits, RBAC, audit logging |
| **Network connectivity issues** | Medium | Low | Implement retry logic, alerting, fallback mechanisms |
| **Secret exposure** | High | Low | Use sealed secrets or external secret managers, regular rotation |
| **Resource exhaustion** | Medium | Medium | Implement resource quotas, monitoring, auto-scaling |

---

## 13. Cost Estimation

### Infrastructure Costs (Monthly, AWS Example)

| Component | Configuration | Estimated Cost |
|-----------|--------------|----------------|
| **EKS Cluster** | 1 cluster | $75 |
| **Worker Nodes** | 3x t3.large (production) | $225 |
| **Worker Nodes** | 2x t3.medium (dev/staging) | $75 |
| **Load Balancer** | 1x ALB | $25 |
| **Container Registry** | ECR (100GB storage) | $10 |
| **RDS** | PostgreSQL instances | $150 |
| **Monitoring** | CloudWatch/Prometheus | $50 |
| **Backups** | S3 storage for Velero | $20 |
| **Total** | | **~$630/month** |

### Tooling Costs

| Tool | License | Cost |
|------|---------|------|
| **ArgoCD** | Open Source | Free |
| **Kubernetes** | Open Source | Free |
| **Kustomize** | Open Source | Free |
| **GitHub Actions** | 3000 min/month included | Free - $100/month |
| **Snyk** | Developer plan | $0 - $99/month |
| **Total** | | **$0 - $200/month** |

### One-Time Costs

| Item | Estimated Cost |
|------|----------------|
| **Training and workshops** | $5,000 |
| **Consulting (optional)** | $10,000 - $30,000 |
| **Migration effort** | 8 weeks of team time |

---

## 14. Success Metrics

### Deployment Metrics
- **Deployment frequency:** Target 10+ deployments per day to dev
- **Lead time for changes:** < 1 hour from commit to dev deployment
- **Mean time to recovery (MTTR):** < 15 minutes with automated rollback
- **Change failure rate:** < 5% of deployments causing issues

### Operational Metrics
- **Service availability:** 99.9% uptime SLA
- **Deployment success rate:** > 95%
- **Automated vs manual deployments:** > 90% automated
- **Time saved per deployment:** ~30 minutes per deployment

### Team Metrics
- **Onboarding time for new developers:** < 1 day
- **Confidence in deployments:** Survey-based improvement
- **Incident response time:** 50% reduction

---

## 15. Recommendations

1. **Start with Development Environment**
   - Deploy to dev first to gain experience
   - Iterate on configuration before moving to production

2. **Implement Progressive Delivery**
   - Start with rolling updates
   - Gradually introduce canary deployments for critical services
   - Use feature flags for additional safety

3. **Invest in Monitoring**
   - Set up comprehensive monitoring before production deployment
   - Create dashboards for all key metrics
   - Implement alerting for critical issues

4. **Documentation is Critical**
   - Document all processes and runbooks
   - Create troubleshooting guides
   - Maintain up-to-date architecture diagrams

5. **Security First**
   - Implement secret management from day one
   - Use RBAC for all access control
   - Regular security audits and scanning

6. **Team Training**
   - Conduct workshops on GitOps, Kubernetes, and ArgoCD
   - Pair programming during initial implementation
   - Create internal knowledge base

---

## 16. Conclusion

Implementing a GitOps-based deployment pipeline with ArgoCD and Kubernetes will significantly improve the MVTA platform's deployment process, reliability, and scalability. The proposed architecture leverages industry best practices and modern tooling to create a robust, automated infrastructure that can support the platform's growth.

**Key Takeaways:**
- GitOps provides a declarative, version-controlled approach to infrastructure
- ArgoCD automates synchronization and provides self-healing capabilities
- Kubernetes offers scalable, resilient container orchestration
- The phased migration approach minimizes risk
- Long-term benefits far outweigh initial implementation costs

**Next Steps:**
1. Review and approve this proposal
2. Allocate resources and budget
3. Begin Phase 1: Infrastructure setup
4. Regular check-ins and progress reviews

---

## Appendix

### A. Glossary

- **GitOps:** Operational framework using Git as single source of truth for declarative infrastructure
- **ArgoCD:** Kubernetes-native continuous delivery tool
- **Kustomize:** Kubernetes native configuration management tool
- **HPA:** Horizontal Pod Autoscaler - automatically scales pods based on metrics
- **RBAC:** Role-Based Access Control
- **Service Mesh:** Infrastructure layer handling service-to-service communication

### B. References

- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kustomize Documentation](https://kustomize.io/)
- [GitOps Principles](https://opengitops.dev/)
- [CNCF Best Practices](https://www.cncf.io/)

### C. Contact Information

**Project Team:**
- DevOps Lead: [Name]
- Platform Architect: [Name]  
- Security Engineer: [Name]

**Stakeholders:**
- Engineering Manager: [Name]
- Product Owner: [Name]
- CTO: [Name]

---

*Document Version: 1.0*  
*Last Updated: February 1, 2026*  
*Next Review: March 1, 2026*
