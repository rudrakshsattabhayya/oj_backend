from django.db import models

# Create your models here.

class CodeModel(models.Model):
    inputs = models.FileField(upload_to="inputs", null = True)
    code = models.FileField(upload_to="codes",  null = True)
    outputs = models.FileField(upload_to="outputs", null = True)
    correctOutputs = models.FileField(upload_to="correctOutputs", null = True)