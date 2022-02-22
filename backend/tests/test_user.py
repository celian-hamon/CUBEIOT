from cgi import test
import json
import requests
from requests.structures import CaseInsensitiveDict
from random import randint

# Endpoints checked :
# POST   - Create user              : user/
# PUT    - Update user              : user/<uid>
# GET    - Get user                 : user/<uid>
# DELETE - Delete user              : user/<uid>


def main():
    base_url = "http://localhost:5000/"
    test_profile = "13"
    test_profile_pd = "admin"

    # CREATE USER
    url = base_url + "user"

    test_user = str(randint(1, 60))
    post_data = """
    {
        "name": "test_user_%s",
        "admin": "0",
        "password": "test_user"
    }
    """ % (test_user)
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("POST", url, post_data=post_data, headers=headers)
    json_body = json.loads(resp.text)
    test_user = str(json_body["id"])
    print("Succefully Created user : " + test_user)

    # UPDATE USER
    url = base_url + "user/" + test_user

    post_data = """
    {
        "name": "test_user_%s_updated",
        "admin": "1"
    }
    """ % (test_user)
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("PUT", url, post_data=post_data, headers=headers)
    json_body = json.loads(resp.text)
    print("Succefully Updated user : " + test_user)

    # GET USER
    url = base_url + "user/" + test_user

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("GET", url, headers=headers)
    json_body = json.loads(resp.text)
    if json_body["name"] != "test_user_%s_updated" % (test_user):
        raise Exception("Error : user not updated")
    print("Succefully Get user : " + test_user)

    # DELETE USER
    url = base_url + "user/" + test_user

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("DELETE", url, headers=headers)
    print("Succefully Delete user : " + test_user)


def call(method, url, post_data=None, headers=CaseInsensitiveDict()):
    if method == "POST":
        resp = requests.post(url, headers=headers, data=post_data)
    if method == "GET":
        resp = requests.get(url, headers=headers)
    if method == "PUT":
        resp = requests.put(url, headers=headers, data=post_data)
    if method == "DELETE":
        resp = requests.delete(url, headers=headers)

    if not resp.status_code or resp.status_code < 200 or resp.status_code > 299:
        raise Exception("Error : " + str(resp.status_code))

    return resp
