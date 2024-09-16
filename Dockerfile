# Use a full-featured Python image instead of alpine
FROM python:3.9-slim

# Install necessary system dependencies
RUN apt-get update && apt-get install -y \
    python3-pip \
    python3-dev \
    libffi-dev \
    libssl-dev \
    gcc \
    make \
    musl-dev \
    && apt-get clean

# Install docker-compose
RUN pip install docker-compose

# Set the working directory inside the container
WORKDIR /app

# Copy the contents of the current directory to the working directory
COPY . /app

EXPOSE 8001

# Command to run docker-compose up
CMD ["docker-compose", "up"]
