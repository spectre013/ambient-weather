from flask import render_template
from flask_sock import Sock
from datetime import datetime
import time
from zoneinfo import ZoneInfo
from forecast import get_forecast
from models import Weather, Stats, Alert
import util

sock = Sock()


@sock.route('/ws')
def echo(ws):
    while True:
        forecast = get_forecast()
        units = "imperial"
        local = ZoneInfo("America/Denver")
        now = datetime.now()
        weather = Weather.query.order_by(Weather.recorded.desc()).first()
        stats = Stats.query.all()
        minmax = util.minmax(stats)
        weather.recorded = weather.recorded.astimezone(local)
        box = util.box_format(units)
        alerts = util.query_list(Alert.query.filter(Alert.ends >= datetime.now()).all())
        wind = util.get_wind()
        lightning_month = util.get_lightning_month()
        btrend = util.trend("baromrelin")
        astro = util.get_astro()

        ws.send(render_template('application/weather.html', now=now,weather=weather,wind=wind,
                                box=box, forecast=forecast, units=units,minmax=minmax, alerts=alerts, 
                                lightning_month=lightning_month, btrend=btrend, astro=astro))
        time.sleep(60)
