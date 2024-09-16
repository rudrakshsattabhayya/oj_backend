from django.shortcuts import render
from rest_framework.views import APIView
from .models import *
from .serializers import *
import jwt
import os
from rest_framework import status
from rest_framework.response import Response
from google.auth import crypt
from google.auth import jwt as gjwt
from uuid import uuid4
import docker
import tarfile
from pathlib import Path
from django.core.files import File
import bcrypt
import requests
import json
from env import (SECRET_KEY, DJANGO_PASSWORD, DJANGO_EVALUATION_SERVER_URL)

BASE_DIR = Path(__file__).resolve().parent.parent

def authenticate(recievedJWT):
    decodedJWT=None
    try:
        decodedJWT = jwt.decode(recievedJWT, SECRET_KEY, algorithms=['HS256'])
    except:
        return {"message": "Login Expired!", "status": status.HTTP_400_BAD_REQUEST}
    
    token = decodedJWT['token']
    user = UserModel.objects.filter(token=token).first()
    if user:
        return {"user": user, "status": status.HTTP_200_OK}
    else:
        return {"message": "User is Invalid!", "status": status.HTTP_404_NOT_FOUND}

def hash_password(password):
    # Generate a salt
    salt = bcrypt.gensalt()

    # Hash the password with the salt
    hashed_password = bcrypt.hashpw(password.encode('utf-8'), salt)

    # Return the hashed password as a string
    return hashed_password.decode('utf-8')

def verify_password(password, hashed_password):
    # Verify the password
    verifiedStatus = False
    try:
        verifiedStatus = bcrypt.checkpw(password.encode('utf-8'), hashed_password.encode('utf-8'))
    except:
        return False
    
    return verifiedStatus

class AuthenticateRoute(APIView):
    def post(self, request):
        recievedToken = request.data['token']
        if not recievedToken:
            return Response({"message": "Token not found! Try to Re-Login.", "status": status.HTTP_404_NOT_FOUND})
        
        res = authenticate(recievedToken)
        if res["status"] == status.HTTP_404_NOT_FOUND:
            return Response(res)
        
        user = res['user']
        obj = {
            "status": res["status"],
            "name": user.name,
            "email": user.email,
            "profilePic": user.profilePic,
            "username": user.username
        }
        return Response(obj)

class AuthenticateRouteForAdmin(APIView):
    def post(self, request):
        recievedToken = request.data['token']
        if not recievedToken:
            return Response({"message": "Token not found! Try to Re-Login.", "status": status.HTTP_404_NOT_FOUND})
        
        res = authenticate(recievedToken)
        user = res['user']
        if res["status"] == status.HTTP_404_NOT_FOUND:
            return Response(res)
        
        if user.isAdmin:
            obj = {
            "status": res["status"],
            "name": user.name,
            "email": user.email,
            "profilePic": user.profilePic
            }
            return Response(obj)
        else:
            return Response({"message": "Requested page can only be accessed by the Admin!", "status": status.HTTP_403_FORBIDDEN})

class LoginView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        googleObj = None
        try:
            googleObj = gjwt.decode(recievedJWT, verify=False)
        except:
            return Response({"message":"Invalid JWT!", "status": status.HTTP_404_NOT_FOUND})
        
        userToken = str(uuid4())
        found = UserModel.objects.filter(email=googleObj['email']).first()
        user = found

        if not found:
            user = UserModel(email=googleObj['email'], name=googleObj['name'], profilePic=googleObj['picture'], username=googleObj['name'], token=userToken)
        else:
            user.token = userToken
        
        user.save()

        payload = {
                    "token": userToken
                  }
        jwtToken = jwt.encode(payload, SECRET_KEY, algorithm='HS256')

        return Response({"jwtToken": jwtToken, "name": user.name, "email": user.email, "profilePic": user.profilePic, "status": status.HTTP_200_OK})

class LoginWithPassword(APIView):
    def post(self, request):
        email = request.data["email"]
        password = request.data["password"]
        user = UserModel.objects.filter(email=email).first()

        if not user:
            return Response({"message": "This email is not registered! Login with google to register yourself.", "status": status.HTTP_404_NOT_FOUND})

        if not user.hashedPassword:
            return Response({"message": "Login with google to add your password!", "status": status.HTTP_501_NOT_IMPLEMENTED})
        
        verified = verify_password(password, user.hashedPassword)

        if verified:
            userToken = str(uuid4())
            user.token = userToken
            user.save()

            payload = {
                    "token": userToken
                    }
            jwtToken = jwt.encode(payload, SECRET_KEY, algorithm='HS256')
            return Response({"jwtToken": jwtToken, "name": user.name, "email": user.email, "profilePic": user.profilePic, "status": status.HTTP_200_OK})
        else:
            return Response({"message": "Password is wrong!", "status": status.HTTP_400_BAD_REQUEST})

class ChangeThePassword(APIView):
    def post(self, request):
        password = request.data["password"]
        recievedJWT = request.data['jwtToken']

        authentication = authenticate(recievedJWT=recievedJWT)
        if authentication['status'] != status.HTTP_200_OK:
            return Response(authentication)
        
        user = authentication["user"]
        hashedPassword = hash_password(password)
        user.hashedPassword = hashedPassword
        user.save()

        return Response({"message": "Password is changed!", "status": status.HTTP_200_OK})        
        


class ChangeUserNameView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)
        if response['status'] == status.HTTP_404_NOT_FOUND:
            return Response({"message" : "User is Invalid!", "status": status.HTTP_404_NOT_FOUND})
        
        newUserName = request.data["newUserName"]

        check = UserModel.objects.filter(username=newUserName).first()
        if check:
            return Response({"message": "A User with this username already exists!", "status": status.HTTP_400_BAD_REQUEST})

        user = UserModel.objects.filter(email=response['user'].email).first()
        user.username = newUserName
        user.save()

        ser_data = UserModelSerializer(user)
        return Response({"username": ser_data.data["username"], "message": "Username successfully changed!", "status": status.HTTP_200_OK})

class CreateProblemView(APIView): 
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] != status.HTTP_200_OK or response['user'].isAdmin == False:
            return Response({"message" : "User is Invalid!", "status": status.HTTP_404_NOT_FOUND})

        data = request.data
        files = request.FILES

        problemStatement = files["problemStatement"]
        hiddenTestCases = files["hiddenTestCases"]
        correctSolution = files["correctSolution"]
        visibleOutputs = files["visibleOutputs"]
        visibleTestCases = files["visibleTestCases"]

        newProblem = ProblemModel(title=data["title"], problemStatement=problemStatement, visibleTestCases=visibleTestCases,
                                  hiddenTestCases=hiddenTestCases, correctSolution=correctSolution,
                                  difficulty=data["difficulty"], visibleOutputs=visibleOutputs)
        newProblem.save()


        emptyFilePath = f"{BASE_DIR}/Uploads/emptyFile.txt"

        with open(emptyFilePath, 'w') as file:
            pass
        
        with open(emptyFilePath, 'rb') as file_obj:
            newProblem.correctOutput.save("correctOutputs.txt", File(file_obj))

        inputs = newProblem.hiddenTestCases
        code = newProblem.correctSolution
        inputs = newProblem.hiddenTestCases
        outputs_path = newProblem.correctOutput.path

        files = {
            "code": code,
            "inputs": inputs,
        }
        obj = {
            "password": DJANGO_PASSWORD,
        }

        respFromEval = requests.post(f"{DJANGO_EVALUATION_SERVER_URL}/get-outputs", files=files, data=obj)
        response_dict = respFromEval.json()
        outputString = response_dict["outputs"]

        with open(outputs_path, 'w') as file:
            file.write(outputString)

        tagslist = data["tags"].split(", ")
        tagsQueryset = TagModel.objects.all()

        for tag in tagslist:
            found = tagsQueryset.filter(name=tag).first()
            y = found
            if(not found):
                y = TagModel(name=tag)
                y.save()
            else:
                y = found

            newProblem.tags.add(y)

        newProblem.save()
        return Response({"QuestionId": newProblem.id, "status": status.HTTP_200_OK})

class GetLeaderBoardView(APIView):
    def get(self, request):
        try:
            users = UserModel.objects.all()
            ser_data = GetLeaderBoardViewSerializer(users, many=True)

            return Response({"response": ser_data.data, "status": status.HTTP_200_OK})
        
        except:
            return Response({"message": "Unable to get the leaderboard!", "status": status.HTTP_400_BAD_REQUEST})

class ListProblemsView(APIView):
    def get(self, request):
        try:
            problems = ProblemModel.objects.all()
            ser_data = ListProblemViewSerializer(problems, many=True)
            
            return Response({"problems": ser_data.data, "status": status.HTTP_200_OK})
        
        except:
            return Response({"message": "Unable to get the Problems!", "status": status.HTTP_400_BAD_REQUEST})
    
class ListSubmissionsView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] != status.HTTP_200_OK:
            return Response(response)
        
        user = response['user']

        submissions = SubmissionModel.objects.filter(user__email = user.email).all()
        ser_data = ListSubmissionsViewSerializer(submissions, many=True)

        return Response({"submissions": ser_data.data, "status": status.HTTP_200_OK})
    
class ListTagsView(APIView):
    def get(self, request):
        try:
            availableTags = TagModel.objects.all()
            ser_data = TagModelSerializer(availableTags, many=True)
            return Response({"tags": ser_data.data, "status": status.HTTP_200_OK})
        except:
            return Response({"message": "Unable to get the Filter Tags!", "status": status.HTTP_400_BAD_REQUEST})
    
class ShowProblemView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] != status.HTTP_200_OK:
            return Response(response)
        
        user = response['user']

        questionId = request.data["questionId"]
        problem = None
        try:
            problem = ProblemModel.objects.filter(id=questionId).first()
        except:
            return Response({"message": "Problem ID is invalid!", "status": status.HTTP_404_NOT_FOUND})

        if(not problem):
            return Response({"message": "Problem ID is invalid!", "status": status.HTTP_404_NOT_FOUND})

        userSubmissions = SubmissionModel.objects.filter(user=user, problem=problem)
        ser_submissions = SubmissionsOfAProblemSerializer(userSubmissions, many=True)

        ser_data = ShowProblemViewSerializer(problem)
        return Response({"response": {"problemsData":ser_data.data, "userSubmissions": ser_submissions.data} , "status": status.HTTP_200_OK})

class ShowProblemSolutionView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] != status.HTTP_200_OK:
            return Response(response)
        
        user = response['user']
        questionId = request.data["questionId"]
        problem = None
        try:
            problem = ProblemModel.objects.filter(id=questionId).first()
        except:
            return Response({"message": "Problem ID is invalid!", "status": status.HTTP_404_NOT_FOUND})

        if(not problem):
            return Response({"message": "Problem ID is invalid!", "status": status.HTTP_404_NOT_FOUND})

        problemIdModelObj = ProblemIdModel.objects.filter(problemId=problem.id, user=user).first()
        if not problemIdModelObj:
            problemIdModelObj = ProblemIdModel(problemId=problem.id, user=user)
            problemIdModelObj.save()

        # correct_solution_url = problem.correctSolution.url if problem.correctSolution else None
        return Response({"solution": problem.correctSolution.url, "status": status.HTTP_200_OK})

class SubmitProblemView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] == status.HTTP_404_NOT_FOUND:
            return Response({"message" : "User is Invalid!", "status": status.HTTP_404_NOT_FOUND})

        submissionFile = request.FILES["code"]
        questionId = request.data["questionId"]
        user = response['user']
        problem = ProblemModel.objects.filter(id=questionId).first()

        if not problem:
            return Response({"message" : "Problem ID is Invalid!", "status": status.HTTP_404_NOT_FOUND})

        submissionObj = SubmissionModel(code=submissionFile, user = user, problem = problem)
        submissionObj.save()

        inputs = problem.hiddenTestCases
        correctOutputs = problem.correctOutput
        code = submissionObj.code

        files = {
            "code": code,
            "inputs": inputs,
            "correctOutputs": correctOutputs
        }
        data = {
            "password": DJANGO_PASSWORD,
        }

        respFromEval = requests.post(f"{DJANGO_EVALUATION_SERVER_URL}/get-verdict", files=files, data=data)
        temp = None
        for x in respFromEval:
            temp = x

        temp2 = temp.decode('utf-8')
        response_dict = json.loads(temp2)
        verdict = response_dict["verdict"]


        problem.totalSubmissions += 1
        user.totalSubmissions += 1

        if verdict:
            solutionViewed = ProblemIdModel.objects.filter(problemId=problem.id, user=user).all()
            submissions = SubmissionModel.objects.filter(user__email = user.email, problem=problem).all()
            if len(solutionViewed) == 0 and len(submissions) == 1:
                proofOfSolved = ProblemIdModel(problemId=problem.id, user=user)
                proofOfSolved.save()
                user.leaderBoardScore = user.leaderBoardScore + problem.difficulty
            
            submissionObj.verdict = True
            problem.acceptedSubmissions += 1
            user.acceptedSubmissions += 1

        user.save()
        problem.save()
        submissionObj.save()

        return Response({"verdict": verdict, "message": "Successful submission!", "status": status.HTTP_200_OK})
    
class DeleteSubmissionsView(APIView):
    def post(self, request):
        recievedJWT = request.data['jwtToken']
        response = authenticate(recievedJWT=recievedJWT)

        if response['status'] != status.HTTP_200_OK or response['user'].isAdmin == False:
            return Response({"message" : "User is Invalid!", "status": status.HTTP_404_NOT_FOUND})
        
        submissionsToBeDeleted = request.data['submissionsToBeDeleted']

        for submissionId in submissionsToBeDeleted:
            submissionsList = SubmissionModel.objects.filter(id=submissionId)
            if not submissionsList.first():
                pass

            submission = submissionsList.first()
            codePath = submission.code.path
            outputsPath = submission.outputs.path

            os.remove(codePath)
            os.remove(outputsPath)

            submission.delete()
        
        return Response({"message":"Selected submissions are deleted!", "status": status.HTTP_200_OK})

def addTags(tags, newProblem):
    tagslist = tags.split(", ")
    tagsQueryset = TagModel.objects.all()

    for tag in tagslist:
        found = tagsQueryset.filter(name=tag).first()
        y = found
        if(not found):
            y = TagModel(name=tag)
            y.save()
        else:
            y = found
        
        newProblem.tags.add(y)
    
    newProblem.save()

def addFilesToProblem(prob, fileName):
    stmt_path = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'stmt.txt')
    with open(stmt_path, 'rb') as stmt_file:
        prob.problemStatement.save('stmt.txt', File(stmt_file))


    visibleTestCasesPath = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'visibleInputs.txt')
    with open(visibleTestCasesPath, 'rb') as visibleTestCasesPathFile:
        prob.visibleTestCases.save('visibleInputs.txt', File(visibleTestCasesPathFile))


    visibleOutputsFilePath = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'visibleOutputs.txt')
    with open(visibleOutputsFilePath, 'rb') as visibleOutputsFile:
        prob.visibleOutputs.save('visibleInputs.txt', File(visibleOutputsFile))


    hiddenTestCasesPath = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'inputs.txt')
    with open(hiddenTestCasesPath, 'rb') as hiddenTestCasesFile:
        prob.hiddenTestCases.save('hiddenTestCases.txt', File(hiddenTestCasesFile))


    correctSolutionPath = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'solution.cpp')
    with open(correctSolutionPath, 'rb') as correctSolutionPathFile:
        prob.correctSolution.save('correctSolution.txt', File(correctSolutionPathFile))

    correctOutputPath = os.path.join(BASE_DIR, 'InitializeProblems', fileName, 'correctOutputs.txt')
    with open(correctOutputPath, 'rb') as correctOutputPathFile:
        prob.correctOutput.save('correctOutput.txt', File(correctOutputPathFile))
    
    prob.save()

def initializeProblem(title, difficulty, tags, fileName, acceptedSubmissions, totalSubmissions):
    prob = ProblemModel(title=title, difficulty=difficulty, acceptedSubmissions=acceptedSubmissions, totalSubmissions=totalSubmissions)
    prob.save()
    
    addTags(tags, prob)
    addFilesToProblem(prob, fileName)

def InitializeAppsData():
    #Create Fake Users for leaderboard
    user = UserModel(email="dummy1@gmail.com", username="Monica", leaderBoardScore=20)
    user.save()
    user = UserModel(email="dummy2@gmail.com", username="Chandler", leaderBoardScore=15)
    user.save()
    user = UserModel(email="dummy3@gmail.com", username="Joey", leaderBoardScore=10)
    user.save()
    user = UserModel(email="dummy4@gmail.com", username="Ross", leaderBoardScore=5)
    user.save()
    user = UserModel(email="dummy5@gmail.com", username="Phoebe", leaderBoardScore=5)
    user.save()
    user = UserModel(email="dummy6@gmail.com", username="Rachel", leaderBoardScore=1)
    user.save()

    #Create Problem Statements
    initializeDirectory = os.path.join(BASE_DIR, 'InitializeProblems')
    if os.path.exists(initializeDirectory):
        subfolders = [folder for folder in os.listdir(initializeDirectory) if os.path.isdir(os.path.join(initializeDirectory, folder))]

        for folder_name in subfolders:
            file_path = os.path.join(BASE_DIR, 'InitializeProblems', folder_name, 'otherinfo.txt')

            with open(file_path, 'r') as file:
                lines = file.readlines()

            info = {}
            for line in lines:
                key, value = line.strip().split('=')
                info[key] = value
            
            title = info.get('name', '')
            tags = info.get('tags', '')
            difficulty = int(info.get('difficulty', 0))
            acceptedSubmissions = int(info.get('acceptedSubmissions', 0))
            totalSubmissions = int(info.get('totalSubmissions', 0))
            initializeProblem(title=title, difficulty=difficulty, tags=tags, fileName=folder_name, acceptedSubmissions=acceptedSubmissions, totalSubmissions=totalSubmissions)
    
    return Response({"message": "Initialized the app data!", "status": status.HTTP_200_OK})

class DontSleep(APIView):
    def get(self, request):
        users = UserModel.objects.all()
        if len(users) >= 6:
            return Response("Thank you for helping me not to sleep!")
        else:
            return InitializeAppsData()
