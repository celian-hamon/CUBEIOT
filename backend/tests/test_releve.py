import json
import requests
from requests.structures import CaseInsensitiveDict
from random import randint

# Endpoints checked :
# POST   - Create releve            : releve/
# GET    - Get releve               : releve/<uid>
# DELETE - Delete releve            : releve/<uid>


def main():
    base_url = "http://localhost:5000/"
    test_profile = "2"
    test_sonde = "1"

    # CREATE RELEVE
    url = base_url + "releve"

    temprature = str(randint(-70, 60))
    humidite = str(randint(1, 100))
    post_data = """
    {
        "temperature":%s,
        "humidite":%s
    }
    """ % (temprature, humidite)
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    headers["profile"] = test_profile
    headers["sonde"] = test_sonde

    resp = call("POST", url, post_data, headers)
    json_body = json.loads(resp.text)
    test_releve = str(json_body["id"])
    print("Succefully Created releve : " + test_releve)

    # GET RELEVE
    url = base_url + "releve/" + test_releve

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile

    resp = call("GET", url, headers=headers)
    json_body = json.loads(resp.text)
    if json_body["temperature"] != temprature and json_body["humidite"] != humidite:
        raise Exception("Error : releve not created correctly")
    else:
        print("Succefully Retreived releve : " + test_releve)

    # DELETE RELEVE
    url = base_url + "releve/" + test_releve

    headers = CaseInsensitiveDict()
    headers["profile"] = test_profile

    resp = call("DELETE", url, headers=headers)
    print("Succefully Deleted releve : " + test_releve)


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
