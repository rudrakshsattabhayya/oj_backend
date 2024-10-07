from .models import *
from pathlib import Path
from django.core.files import File
import subprocess
import os
from rest_framework import status
from celery import shared_task
import requests
from env import (DJANGO_BACKEND_SERVER_URL)
from django.conf import settings
import resource
import shutil

BASE_DIR = Path(__file__).resolve().parent.parent


@shared_task
def evalutate(codeModalId, time_limit = 3, output_size_limit = 50 * 1024 * 1024):
    reason = ""
    verdict = False
    try:
      codeModalObj = CodeModel.objects.get(id=codeModalId)

      cppCodeFilePath, inputsPath, outputsPath, correctOutputsPath = getFilePaths(codeModalObj)

      isCompiled, reason = compileCode(cppCodeFilePath)

      if isCompiled:
        verdict, reason = executeCode(cppCodeFilePath, inputsPath, outputsPath, correctOutputsPath, time_limit, output_size_limit)

    except Exception as e:
        return {"error": str(e), "status": status.HTTP_400_BAD_REQUEST}
    finally:
      cleanup([cppCodeFilePath, inputsPath, outputsPath, correctOutputsPath, f"{cppCodeFilePath}.exe"], codeModalObj)
      UpdateVerdict(task_id=evalutate.request.id, verdict=verdict, reason=reason)

    return {"task_id": evalutate.request.id, "verdict": verdict, "reason": reason}

def getFilePaths(codeModalObj):
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

  return cppCodeFilePath, inputsPath, outputsPath, correctOutputsPath


def compileCode(cppCodeFilePath):
  executableFilePath = f"{cppCodeFilePath}.exe"
  compile_command = f'g++ "{cppCodeFilePath}" -o "{executableFilePath}"'
  
  # Compile the C++ code
  compile_process = subprocess.Popen(compile_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
  compile_output, compile_error = compile_process.communicate()

  if compile_process.returncode != 0:
    reason =  compile_error.decode("utf-8")
    return False, reason
  
  return True, ""

def set_memory_limit(limit):
    def limit_memory():
        resource.setrlimit(resource.RLIMIT_AS, (limit, resource.RLIM_INFINITY))
    return limit_memory

def executeCode(cppCodeFilePath, inputsPath, outputsPath, correctOutputsPath, time_limit, output_size_limit):
    executableFilePath = f"{cppCodeFilePath}.exe"
    execute_command = f'"{executableFilePath}" < "{inputsPath}"'
    verdict = False
    reason = ""

    try:
        # Use pre-exec to limit the memory usage of the child process before execution
        execute_process = subprocess.Popen(
            execute_command,
            preexec_fn=set_memory_limit(output_size_limit),
            shell=True, 
            stdout=subprocess.PIPE, 
            stderr=subprocess.PIPE, 
        )

        execute_output, execute_error = execute_process.communicate(timeout=time_limit)

        # Check if there was an execution error
        if execute_process.returncode != 0:
            return False, execute_error.decode("utf-8")  # Early exit on error

        # Check the size of the output before writing to the file
        output_size = len(execute_output)
        if output_size > output_size_limit:
            return False, "Memory Limit Exceeded!"  # Early exit on output size limit

        # Write the output to a file if it's within the limit
        with open(outputsPath, 'wb') as file:
            file.write(execute_output)

        # Compare the output file with the correct output
        verdict = compareFiles(outputsPath, correctOutputsPath)
        return verdict, reason
    except subprocess.TimeoutExpired:
        return False, "Time Limit Exceeded!"
    except OSError as e:
        if e.errno == 12:
          return False, "Memory Limit Exceeded!"
        return False, e
    except Exception as e:
      return False, e


def UpdateVerdict(task_id, verdict, reason=""):
    response = {
        "task_id": task_id,
        "verdict": verdict,
        "reason": reason
    }

    # print("Task: ", response)

    requests.post(f"{DJANGO_BACKEND_SERVER_URL}/update_verdict", data=response)

def compareFiles(file1_path, file2_path):
    with open(file1_path, 'r') as file1:
        content1 = file1.read()

    with open(file2_path, 'r') as file2:
        content2 = file2.read()

    return content1 == content2

def cleanup(files, codeModalObj):
  folder = f"{BASE_DIR}/Uploads"

  # Loop over all files and directories in the folder
  for filename in os.listdir(folder):
      file_path = os.path.join(folder, filename)

      try:
          if os.path.isfile(file_path) or os.path.islink(file_path):
              # Remove files and symlinks
              os.unlink(file_path)
          elif os.path.isdir(file_path):
              # Remove directories
              shutil.rmtree(file_path)
      except Exception as e:
          print(f'Failed to delete {file_path}. Reason: {e}')
  
  codeModalObj.delete()