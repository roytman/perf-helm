FROM ubuntu
COPY perf-helm /

CMD exec /bin/bash -c "trap : TERM INT; sleep infinity & wait"