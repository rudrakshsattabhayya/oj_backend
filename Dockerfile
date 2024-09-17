# Use Ubuntu as the base image
FROM ubuntu:20.04

# Set non-interactive mode to avoid prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages and Docker
RUN apt-get update && \
    apt-get install -y \
    curl \
    apt-transport-https \
    ca-certificates \
    gnupg-agent \
    software-properties-common \
    && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable" \
    && apt-get update \
    && apt-get install -y docker.io \
    && curl -L "https://github.com/docker/compose/releases/download/$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep tag_name | cut -d\" -f4)/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy your application files into the container
COPY . .

# Command to run docker-compose when the container starts
CMD ["docker-compose", "up"]
