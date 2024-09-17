FROM python:3.11.4

WORKDIR /app

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt

# COPY ../.env .env

COPY . .

EXPOSE 8001

CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:8001"]