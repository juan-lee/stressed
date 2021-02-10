# Stressed
A simple operator that wraps the [stress-ng](https://wiki.ubuntu.com/Kernel/Reference/stress-ng)
tool for running stress on a kubernetes cluster.

## Quickstart

``` bash
export IMG=<your registry>/controller:latest
export STRESS_IMG=<your registry>/stress-ng:latest
export KUBECONFIG=<your kubeconfig path>

# build and push image
make docker-build docker-push

# install and deploy
make install deploy

# deploy the sample
kubectl apply -f config/samples/test_v1alpha1_stress.yaml
```
