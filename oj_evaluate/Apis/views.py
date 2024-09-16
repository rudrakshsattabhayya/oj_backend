from rest_framework.views import APIView
from django.http import JsonResponse
from .models import *
from pathlib import Path
from django.core.files import File
import subprocess
import os
from rest_framework import status
from env import (DJANGO_PASSWORD)

BASE_DIR = Path(__file__).resolve().parent.parent

def compareFiles(file1_path, file2_path):
    with open(file1_path, 'r') as file1:
        content1 = file1.read()

    with open(file2_path, 'r') as file2:
        content2 = file2.read()

    return content1 == content2

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
        code = request.FILES['code']
        inputs = request.FILES['inputs']
        correctOutputs = request.FILES['correctOutputs']
        password = request.data['password']
        
        if not password == DJANGO_PASSWORD:
            return JsonResponse({"message": "Invalid User!", "status": status.HTTP_401_UNAUTHORIZED})

        codeModalObj = CodeModel(code=code, inputs=inputs, correctOutputs=correctOutputs)
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
        correctOutputsPath = codeModalObj.correctOutputs.path

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

        verdict = compareFiles(outputsPath, correctOutputsPath)
        
        os.remove(cppCodeFilePath)
        os.remove(inputsPath)
        os.remove(outputsPath)
        os.remove(correctOutputsPath)
        os.remove(executableFilePath)
        codeModalObj.delete()

        return JsonResponse({"verdict":verdict, "status": status.HTTP_200_OK})      
    
class Heartbeat(APIView):
    def get(self, request):
        return JsonResponse({"status": status.HTTP_200_OK})
