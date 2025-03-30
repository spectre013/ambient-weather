from flask import Flask
from config import config
from db import db
from dotenv import load_dotenv
from .routes import application as application_blueprint
import util
import inspect
import os
from cache import cache

def create_app(environment):
    print("Loading environment", environment)
    if environment == 'development':
        load_dotenv()

    app_config = config[environment]
    print(os.environ.get('DATABASE_URL'))
    app_config.SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL')
    app = Flask(__name__,
                template_folder=app_config.TEMPLATES,
                static_folder=app_config.STATIC)
    app.config.from_object(app_config)
    cache.init_app(app)
    db.init_app(app)
    from ws import sock
    sock.init_app(app)
    app.static_folder = app_config.STATIC

    with app.app_context():
        app.register_blueprint(application_blueprint, url_prefix='/')

        
        for name, func in inspect.getmembers(util, inspect.isfunction):
            app.add_template_global(func, name)


    return app
