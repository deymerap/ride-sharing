#!/bin/bash
SOCKET="${XDG_RUNTIME_DIR}/podman/podman.sock"
if [ ! -S "$SOCKET" ]; then
  echo "ERROR: Podman socket no encontrado en $SOCKET" >&2
  echo "¿Está corriendo Podman?" >&2
  exit 1
fi
export DOCKER_HOST="unix://${SOCKET}"
export DOCKER_BUILDKIT=0
echo "Usando Podman: ${DOCKER_HOST}"
exec tilt up "$@"
