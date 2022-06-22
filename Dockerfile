FROM ubuntu
COPY perf-helm /root
COPY afm /opt/fybrik/charts/ghcr.io/fybrik/arrow-flight-module-chart:0.7.0

CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"