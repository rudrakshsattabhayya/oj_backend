# Use the official Python base image with Django installed
FROM python:3.9

# Set environment variables
ENV PYTHONUNBUFFERED 1

# Set the working directory in the container
WORKDIR /app

# Copy the requirements.txt file to the working directory
COPY requirements.txt .

# Install the required packages for Django
RUN pip install --no-cache-dir -r requirements.txt

# Install the required packages for executing the CPP file
RUN apt-get update && apt-get install -y g++

# Copy the Django project files to the working directory
COPY . .

EXPOSE 8002

# Run Django's migration and then the development server
CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:8002"]

#Running the follwing cmd to ensure security and limiting the resources of container
#docker run --cpus=0.5 --memory=512m --read-only oj_evaluation_server