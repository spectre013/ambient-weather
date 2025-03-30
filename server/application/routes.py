from flask import render_template, Blueprint, jsonify, send_from_directory
from datetime import datetime
import os
import util

application = Blueprint('application', __name__)


@application.route('/favicon.ico')
def favicon():
    return send_from_directory(os.path.join(application.root_path, 'static').replace("application/", ""),
                               'favicon.ico', mimetype='image/vnd.microsoft.icon')


@application.route('/')
def index():
    print("index")
    now = datetime.now()
    return render_template('application/index.html', now=now)


@application.route('/almanac/<sensor>/<timeframe>')
def almanac(sensor, timeframe):
    now = datetime.now()
    chart = util.chart_format(timeframe,sensor)
    return render_template('application/almanac.html', now=now, sensor=sensor, timeframe=timeframe, chart=chart)
