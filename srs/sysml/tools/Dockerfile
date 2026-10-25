# SysML v2 Pilot Implementation: воспроизводимый экспорт настоящих моделей.
FROM condaforge/miniforge3:24.11.3-0
ENV CONDA_REMOTE_READ_TIMEOUT_SECS=300 CONDA_REMOTE_CONNECT_TIMEOUT_SECS=30 CONDA_REMOTE_MAX_RETRIES=5
RUN --mount=type=cache,target=/opt/conda/pkgs conda install -y -c conda-forge jupyter-sysml-kernel=0.62.0 graphviz=2.* nbformat nbclient cairosvg
COPY --from=eclipse-temurin:21.0.4_7-jre /opt/java/openjdk /opt/jdk21
RUN apt-get update && apt-get install -y --no-install-recommends fontconfig fonts-dejavu-core \
    && rm -rf /var/lib/apt/lists/*
ENV JAVA_HOME=/opt/jdk21 PATH=/opt/jdk21/bin:$PATH JAVA_TOOL_OPTIONS=-Djava.awt.headless=true
WORKDIR /srs/sysml
