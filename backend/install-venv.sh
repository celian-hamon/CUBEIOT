#!/bin/bash
python3 -m venv venv

FLASK_APP=main
export FLASK_APP

./venv/bin/pip install -r requirements.txt
