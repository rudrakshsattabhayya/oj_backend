#!/bin/sh

# Start the Celery worker in the background
celery -A oj_evaluation_server worker -B --loglevel=info &

# Run Django's migration and then the development server
python manage.py migrate
python manage.py migrate django_celery_results
python manage.py migrate django_celery_beat
python manage.py runserver 0.0.0.0:8002
