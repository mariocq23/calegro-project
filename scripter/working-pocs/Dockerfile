# Use a minimal Ubuntu image as the base
FROM ubuntu:latest

# Set environment variable to avoid interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# Update package lists and install necessary tools
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    curl \
    gnupg \
    erlang-base \
    erlang-crypto \
    erlang-mnesia \
    erlang-runtime-tools \
    erlang-ssl \
    erlang-syntax-tools \
    erlang-xmerl \
    dosbox \
    redis-server

# Add RabbitMQ repository key and source list
#RUN curl -fsSL https://www.rabbitmq.com/rabbitmq-release-signing-key.asc | gpg --dearmor > /usr/share/keyrings/rabbitmq-keyring.gpg && \
#    echo "deb [signed-by=/usr/share/keyrings/rabbitmq-keyring.gpg] https://packagecloud.io/rabbitmq/rabbitmq-server/ubuntu/ jammy main\" | tee /etc/apt/sources.list.d/rabbitmq.list

# Update again to find the RabbitMQ package and install it
#RUN apt-get update && apt-get install -y rabbitmq-server

# Define a simple script to start services and then launch a shell
RUN echo '#!/bin/bash' > /usr/local/bin/start-env.sh && \
    echo 'echo "Starting Redis..."' >> /usr/local/bin/start-env.sh && \
    echo 'redis-server --daemonize yes' >> /usr/local/bin/start-env.sh && \
    echo 'echo "Starting RabbitMQ..."' >> /usr/local/bin/start-env.sh && \
    #echo 'rabbitmq-server -detached' >> /usr/local/bin/start-env.sh && \
    echo 'echo "Services started. Entering shell..."' >> /usr/local/bin/start-env.sh && \
    echo 'exec bash' >> /usr/local/bin/start-env.sh && \
    chmod +x /usr/local/bin/start-env.sh

# The entrypoint for the container, which runs the startup script
ENTRYPOINT ["/usr/local/bin/start-env.sh"]