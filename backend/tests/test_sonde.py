import json
import requests
from requests.structures import CaseInsensitiveDict
from random import randint

# Endpoints checked :
# POST   - Create sonde              : sonde/
# GET    - Get sonde                 : sonde/<uid>
# PUT    - Update sonde              : sonde/<uid>
# PUT    - Desactivate sonde         : sonde/<uid>/desactivate
# PUT    - Activate sonde            : sonde/<uid>/activate
# DELETE - Delete sonde              : sonde/<uid>


def main():
    base_url = "http://localhost:5000/"
    test_profile = "13"
    test_profile_pd = "admin"

    # CREATE SONDE
    url = base_url + "sonde"

    post_data = """
    {
        "latitude": "123",
        "longitude": "123"
    }
    """
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("POST", url, post_data, headers)
    json_body = json.loads(resp.text)
    test_sonde = str(json_body["id"])
    print("Succefully Created sonde : " + test_sonde)

    # UPDATE SONDE
    url = base_url + "sonde/" + test_sonde

    latitude = str(randint(1, 60))
    longitude = str(randint(1, 60))
    post_data = """
    {
        "latitude": %s,
        "longitude": %s
    }
    """ % (latitude, longitude)
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("PUT", url, post_data, headers)
    json_body = json.loads(resp.text)
    print("Succefully Updated sonde : " + test_sonde)

    # GET SONDE
    url = base_url + "sonde/" + test_sonde

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("GET", url, headers=headers)
    json_body = json.loads(resp.text)
    if json_body["latitude"] != float(latitude) or json_body["longitude"] != float(longitude):
        raise Exception("Error : sonde not updated")
    print("Succefully Retrieved sonde : " + test_sonde)

    # DESACTIVATE SONDE
    url = base_url + "sonde/" + test_sonde + "/desactivate"

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("PUT", url, headers=headers)
    json_body = json.loads(resp.text)
    if not json_body["success"]:
        raise Exception("Error : sonde not desactivated")
    print("Succefully Desactivated sonde : " + test_sonde)

    # ACTIVATE SONDE
    url = base_url + "sonde/" + test_sonde + "/activate"

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile
    headers["password"] = test_profile_pd

    resp = call("PUT", url, headers=headers)
    json_body = json.loads(resp.text)
    if not json_body["success"]:
        raise Exception("Error : sonde not activated")
    print("Succefully Activated sonde : " + test_sonde)

    # DELETE SONDE
    url = base_url + "sonde/" + test_sonde

    resp = call("DELETE", url, headers=headers)
    print("Succefully Deleted sonde : " + test_sonde)


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
