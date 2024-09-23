from django.urls import path
from .views import *

urlpatterns = [
    path('get-verdict', GetTheVerdict.as_view()),
    path('get-outputs', GetTheOutputs.as_view()),
    path('heartbeat', Heartbeat.as_view()),
    path('create_superuser', CreateSuperUser.as_view()),
]