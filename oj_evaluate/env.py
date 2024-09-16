import os
from dotenv import load_dotenv

load_dotenv()

DJANGO_SECRET_KEY = os.getenv("EVALUATION_DJANGO_SECRET_KEY")
DJANGO_PASSWORD = os.getenv("EVALUATION_DJANGO_PASSWORD")
