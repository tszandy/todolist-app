ARG BASE_IMAGE=ghcr.io/devcontainers/go:1-ubuntu
FROM ${BASE_IMAGE}

ARG USERNAME=devuser
ARG USER_UID=1000
ARG USER_GID=1000

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    ca-certificates curl git sudo procps \
 && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid ${USER_GID} ${USERNAME} || true \
 && useradd --uid ${USER_UID} --gid ${USER_GID} -m ${USERNAME} \
 && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
 && chmod 0440 /etc/sudoers.d/${USERNAME}

WORKDIR /workspace
USER ${USERNAME}
