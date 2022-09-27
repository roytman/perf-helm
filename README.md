# perf-helm

This utility is intended to measure [fybrik](https://github.com/fybrik/fybrik) helmer performance.

You can run it locally, see `perf-helm -h` 
In order to run it in a K8s cluster, just execute 
```bash
kubectl apply -f pod.yaml
```
, which will create a pod and a RoleBinding to the default ServiceAccount. 
Login into the pod 
```bash
kubectl exec -it perf-helm -- bash
``` 
and exec 
```bash
/root/perf-helm
``` 

Meantime, I have not resolved a certificate issue when the program is running in a Pod and trying to pull helm charts from a remote repository.
Therefore, I use locally preloaded charts. If you want to use your chart, just store them into `/opt/fybrik/charts/`
