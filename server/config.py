import os

class Config:
    BASE_DIR = os.path.abspath(os.path.dirname(__file__))
    STATIC = os.path.join(os.path.dirname(__file__), "static")
    TEMPLATES = os.path.join(os.path.dirname(__file__), "templates")
    SQLALCHEMY_DATABASE_URI = ""
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    SECRET_KEY = 'iE9S3jq$uASH%MWk'
    # SQLALCHEMY_ECHO = True

config = {
    'production': Config,
    'development': Config,
}
