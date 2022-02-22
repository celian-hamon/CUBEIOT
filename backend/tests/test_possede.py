import json
from typing import Literal
import requests
from requests.structures import CaseInsensitiveDict


# Endpoints checked :
# GET   - Get list of single sonde          : sonde/solo
# DELETE - Delete user from sonde           : sonde/<sonde_uid>/user/<user_uid>
# PUT   - Add user to sonde                 : sonde/<sonde_uid>/user/<user_uid>
# TODO GET   - Get list of users of sonde        : sonde/<sonde_uid>/user
# GET  - Get list of sondes of user         : user/<user_uid>/sonde


def main():
    base_url = "http://localhost:5000/"
    test_profile = "13"
    test_sonde = "1"
    test_profile_pd = "admin"

    # GET LIST OF FATHERLESS SONDE
    url = base_url + "sonde/solo"

    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("GET", url, headers=headers)
    json_body = json.loads(resp.text)
    print("Succefully Retrieved FatherLess Sonde : " + str(json_body))

    # GET LIST OF SONDE OF USER
    url = base_url + "user/" + test_profile + "/sonde"

    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("GET", url, headers=headers)
    json_body = json.loads(resp.text)
    print("Succefully Retrieved Sonde list from user : " + str(json_body))

    is_in_json = is_in_json_body(json_body, test_sonde)
    print(is_in_json)
    if is_in_json:
        # DELETE USER FROM SONDE
        print("USER IS IN SONDE LIST")
        url = base_url + "sonde/" + test_sonde + "/user/" + test_profile

        headers = CaseInsensitiveDict()
        headers["Content-Type"] = "application/json"
        headers["profile"] = test_profile
        headers["password"] = test_profile_pd

        resp = call("DELETE", url, headers=headers)
        print("Succefully Deleted User from Sonde : " + test_sonde)

        # ADD USER TO SONDE
        url = base_url + "sonde/" + test_sonde + "/user/" + test_profile

        headers = CaseInsensitiveDict()
        headers["Content-Type"] = "application/json"
        headers["profile"] = test_profile
        headers["password"] = test_profile_pd

        resp = call("PUT", url, headers=headers)
        json_body = json.loads(resp.text)
        print("Succefully Added user : " + str(test_profile) +
              " to sonde : " + str(test_sonde))

    elif is_in_json == False:
        # ADD USER TO SONDE
        url = base_url + "sonde/" + test_sonde + "/user/" + test_profile

        headers = CaseInsensitiveDict()
        headers["Content-Type"] = "application/json"
        headers["profile"] = test_profile
        headers["password"] = test_profile_pd

        resp = call("PUT", url, headers=headers)
        json_body = json.loads(resp.text)
        print("Succefully Added user : " + str(test_profile) +
              " to sonde : " + str(test_sonde))

        # DELETE USER FROM SONDE
        url = base_url + "sonde/" + test_sonde + "/user/" + test_profile

        headers = CaseInsensitiveDict()
        headers["Content-Type"] = "application/json"
        headers["profile"] = test_profile
        headers["password"] = test_profile_pd

        resp = call("DELETE", url, headers=headers)
        print("Succefully Deleted User from Sonde : " + test_sonde)

    # GET LIST OF USERS OF SONDE


def call(method, url, headers: CaseInsensitiveDict = CaseInsensitiveDict(), post_data: Literal = None) -> requests.Response:
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


def is_in_json_body(json_body, test_sonde):
    for i in range(len(json_body)):
        if str(test_sonde) == str(json_body[i]["id"]):
            return True
    return False
