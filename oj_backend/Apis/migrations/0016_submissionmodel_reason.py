# Generated by Django 3.2.19 on 2024-10-06 15:16

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('Apis', '0015_auto_20240922_1451'),
    ]

    operations = [
        migrations.AddField(
            model_name='submissionmodel',
            name='reason',
            field=models.TextField(null=True),
        ),
    ]
