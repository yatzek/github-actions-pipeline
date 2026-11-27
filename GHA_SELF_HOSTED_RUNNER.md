# Self-hosted Github actions runners in minikube

## Install ARC

```bash
helm upgrade --install arc \
  --namespace arc-systems \
  --create-namespace \
  --version 0.15.1 \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller
```

```
App ID: 5952209

https://github.com/organizations/yatzek/settings/installations/5952209

```

## Install autoscaling runner set in DIND mode

```bash
helm upgrade --install arc-autoscaled-runners \
  --namespace arc-systems \
  --create-namespace \
  -f values.yaml \
  --version 0.15.1 \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set
```

## Docs

- [actions runner controller docs](https://docs.github.com/en/actions/tutorials/use-actions-runner-controller/get-started)
- [authentication to github api](https://docs.github.com/en/actions/how-tos/manage-runners/use-actions-runner-controller/authenticate-to-the-api)
- [charts](https://github.com/actions/actions-runner-controller/tree/master/charts)
