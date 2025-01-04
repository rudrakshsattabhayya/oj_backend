from rest_framework.views import APIView
from django.http import JsonResponse
from .models import *
from .evaluate import evalutate
from pathlib import Path
from django.core.files import File
import subprocess
import os
from rest_framework import status
from env import (DJANGO_PASSWORD)
from django.contrib.auth.models import User
from django.core.files.uploadedfile import InMemoryUploadedFile
from django.conf import settings
import requests
from django.core.files.base import ContentFile
from django.http import JsonResponse
from rest_framework import status

BASE_DIR = Path(__file__).resolve().parent.parent

class GetTheOutputs(APIView):
    def post(self, request):
        code = request.FILES['code']
        inputs = request.FILES['inputs']
        password = request.data['password']
        
        if not password == DJANGO_PASSWORD:
            return JsonResponse({"message": "Invalid User!", "status": status.HTTP_401_UNAUTHORIZED})

        codeModalObj = CodeModel(code=code, inputs=inputs)
        codeModalObj.save()

        emptyFilePath = f"{BASE_DIR}/Uploads/emptyFile.txt"

        with open(emptyFilePath, 'w') as file:
            pass

        with open(emptyFilePath, 'rb') as file_obj:
            codeModalObj.outputs.save("new_file.txt", File(file_obj))

        codeModalObj.save()

        cppCodeFilePath = codeModalObj.code.path
        inputsPath = codeModalObj.inputs.path
        outputsPath = codeModalObj.outputs.path

        executableFilePath = f"{cppCodeFilePath}.exe"
        compile_command = f'g++ "{cppCodeFilePath}" -o "{executableFilePath}"'  # Updated compile_command
        compile_process = subprocess.Popen(compile_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        compile_output, compile_error = compile_process.communicate()

        if compile_process.returncode != 0:
            return JsonResponse({"message": compile_error.decode("utf-8"), "status": status.HTTP_400_BAD_REQUEST})

        execute_command = f'"{executableFilePath}" < "{inputsPath}"'
        execute_process = subprocess.Popen(execute_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        execute_output, execute_error = execute_process.communicate()

        if execute_process.returncode != 0:
            return JsonResponse({"message": execute_error.decode("utf-8"), "status": status.HTTP_400_BAD_REQUEST})
        
        with open(outputsPath, 'wb') as file:
            file.write(execute_output)
        
        output_string=None
        with open(outputsPath, 'r') as file:
            output_string = file.read()

        os.remove(cppCodeFilePath)
        os.remove(inputsPath)
        os.remove(outputsPath)
        os.remove(executableFilePath)
        codeModalObj.delete()

        return JsonResponse({"outputs":output_string, "status": status.HTTP_200_OK})

class GetTheVerdict(APIView):
    def post(self, request):
        if not request.data['password'] == DJANGO_PASSWORD:
            return {"message": "Invalid User!", "status": status.HTTP_401_UNAUTHORIZED}

        try:
            def get_file_from_url(url):
                """Helper function to fetch file content from a URL."""
                try:
                    # Download the file from URL
                    response = requests.get(url)
                    response.raise_for_status()  # Raises an HTTPError for bad responses
                    
                    # Get the filename from the URL
                    filename = os.path.basename(url)
                    if not filename:
                        filename = 'downloaded_file'  # Default filename if none found in URL
                    
                    # Create a Django ContentFile
                    file_content = ContentFile(response.content)
                    file_content.name = filename
                    
                    return file_content
                except requests.exceptions.RequestException as e:
                    raise Exception(f"Error downloading file from URL: {str(e)}")
            
            # Check if the field is a URL or a file upload
            def process_file_field(field_name):
                field_value = request.data.get(field_name)

                if isinstance(field_value, str) and field_value.startswith('http'):
                    return get_file_from_url(field_value)
                else:
                    return field_value

            code_file = process_file_field('code')
            inputs_file = process_file_field('inputs')
            correct_outputs_file = process_file_field('correctOutputs')

            codeModalObj = CodeModel(
                code=code_file,
                inputs=inputs_file,
                correctOutputs=correct_outputs_file
            )
            codeModalObj.save()

            task = evalutate.delay(codeModalObj.id)
            return JsonResponse({"task_id": task.id, "status": status.HTTP_200_OK})

        except Exception as e:
            return JsonResponse({"error": str(e)}, status=status.HTTP_400_BAD_REQUEST)

class CreateSuperUser(APIView):
    def post(self, request):
        django_pwd = request.data['django_pwd']
        
        if not django_pwd == DJANGO_PASSWORD:
            return JsonResponse({"message": "Invalid User!", "status": status.HTTP_401_UNAUTHORIZED})
        
        username = request.data['username']
        email = request.data['email']
        password = request.data['password']

        User.objects.create_superuser(username, email, password)

        return JsonResponse({"status": status.HTTP_200_OK})
    
class Heartbeat(APIView):
    def get(self, request):
        return JsonResponse({"status": status.HTTP_200_OK})
