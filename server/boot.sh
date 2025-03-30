#!/bin/bash
gunicorn --workers=2 --bind 0.0.0.0:3000 'application:create_app("production")'
