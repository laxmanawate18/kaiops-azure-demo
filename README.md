# kaiops-azure-demo

A tiny stateless Go HTTP API used to exercise the KaiOps deployment pipeline on
**Azure (AKS)**:

```
GitHub repo → CI (build image → Artifact Registry → push-back manifest tag)
→ ArgoCD Application (destination: AKS cluster my-demo-cluster)
→ AKS namespace kaiops-azure-demo → KaiOps poller (cloud_provider=azure)
→ RCA → Slack thread → HITL approve → Azure executor (AKS) remediates
```

## Endpoints

- `GET /` — `{"service":"kaiops-azure-demo","status":"ok"}`
- `GET /healthz` — `{"status":"healthy","app":"kaiops-azure-demo"}`
- `GET /api/info` — app name, version, pod, host, timestamp

## Deploy

- Manifest: `k8s/deployment.yaml` (Namespace `kaiops-azure-demo` + Deployment + Service)
- ArgoCD Application: `argocd-application.yaml` (dest: AKS `my-demo-cluster`)
- CI: `.github/workflows/build-deploy.yml` (authenticates to GCP Artifact Registry
  via Workload Identity Federation, then pushes back the image tag for ArgoCD sync)
