# perf-helm

This utility is intended to measure [fybrik](https://github.com/fybrik/fybrik) helmer performance.

You can run it locally, see `perf-helm -h` 
In order to run it in a K8s cluster, just execute `kubectl -f pod.yam`, which will create a pod with existing binaries 
and the test chart.
Login into the pod `kubectll exec -it perf-helm -- bash` and exec `/root/perf-helm` 

Meantime, I have not resolved a certificate issue when the program is running in a Pod and trying to pull helm charts from a repository.
Therefore, I use locally preloaded charts. If you want to use your chart, just store it into `/opt/fybrik/charts/`