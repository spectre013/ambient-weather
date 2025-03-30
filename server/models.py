from datetime import datetime, timedelta
from db import db


class Weather(db.Model):
    __tablename__ = 'records'
    id = db.Column(db.Integer, primary_key=True)
    mac = db.Column(db.String, nullable=False)
    recorded = db.Column(db.DateTime(timezone=True), nullable=False, default=datetime.utcnow)
    baromabsin = db.Column(db.Float)
    baromrelin = db.Column(db.Float)
    battout = db.Column(db.Integer)
    batt1 = db.Column(db.Integer)
    batt2 = db.Column(db.Integer)
    batt3 = db.Column(db.Integer)
    batt4 = db.Column(db.Integer)
    batt5 = db.Column(db.Integer)
    batt6 = db.Column(db.Integer)
    batt7 = db.Column(db.Integer)
    batt8 = db.Column(db.Integer)
    batt9 = db.Column(db.Integer)
    batt10 = db.Column(db.Integer)
    co2 = db.Column(db.Float)
    dailyrainin = db.Column(db.Float)
    dewpoint = db.Column(db.Float)
    eventrainin = db.Column(db.Float)
    feelslike = db.Column(db.Float)
    hourlyrainin = db.Column(db.Float)
    hourlyrain = db.Column(db.Float)
    humidity = db.Column(db.Integer)
    humidity1 = db.Column(db.Integer)
    humidity2 = db.Column(db.Integer)
    humidity3 = db.Column(db.Integer)
    humidity4 = db.Column(db.Integer)
    humidity5 = db.Column(db.Integer)
    humidity6 = db.Column(db.Integer)
    humidity7 = db.Column(db.Integer)
    humidity8 = db.Column(db.Integer)
    humidity9 = db.Column(db.Integer)
    humidity10 = db.Column(db.Integer)
    humidityin = db.Column(db.Integer)
    lastrain = db.Column(db.DateTime)
    maxdailygust = db.Column(db.Float)
    relay1 = db.Column(db.Integer)
    relay2 = db.Column(db.Integer)
    relay3 = db.Column(db.Integer)
    relay4 = db.Column(db.Integer)
    relay5 = db.Column(db.Integer)
    relay6 = db.Column(db.Integer)
    relay7 = db.Column(db.Integer)
    relay8 = db.Column(db.Integer)
    relay9 = db.Column(db.Integer)
    relay10 = db.Column(db.Integer)
    monthlyrainin = db.Column(db.Float)
    solarradiation = db.Column(db.Float)
    tempf = db.Column(db.Float)
    temp1f = db.Column(db.Float)
    temp2f = db.Column(db.Float)
    temp3f = db.Column(db.Float)
    temp4f = db.Column(db.Float)
    temp5f = db.Column(db.Float)
    temp6f = db.Column(db.Float)
    temp7f = db.Column(db.Float)
    temp8f = db.Column(db.Float)
    temp9f = db.Column(db.Float)
    temp10f = db.Column(db.Float)
    tempinf = db.Column(db.Float)
    totalrainin = db.Column(db.Float)
    uv = db.Column(db.Float)
    weeklyrainin = db.Column(db.Float)
    winddir = db.Column(db.Integer)
    windgustmph = db.Column(db.Float)
    windgustdir = db.Column(db.Integer)
    windspeedmph = db.Column(db.Float)
    yearlyrainin = db.Column(db.Float)
    battlightning = db.Column(db.Integer)
    lightningday = db.Column(db.Integer)
    lightninghour = db.Column(db.Integer)
    lightningtime = db.Column(db.DateTime)
    lightningdistance = db.Column(db.Float)
    aqipm25 = db.Column(db.Integer)
    aqipm2524h = db.Column(db.Integer)

    def __repr__(self):
        return f"<Records {self.id} - {self.recorded}>"

class Stats(db.Model):
    id = db.Column(db.String, primary_key=True)  # character varying
    recorded = db.Column(db.DateTime, nullable=False, default=datetime.utcnow)  # timestamp without time zone
    value = db.Column(db.Numeric, nullable=False)  # numeric

    def __repr__(self):
        return f"<Record {self.id} - {self.recorded} - {self.value}>"

class Alert(db.Model):
    __tablename__ = 'alerts'
    
    id = db.Column(db.Integer, primary_key=True)
    alertid = db.Column(db.String)
    wxtype = db.Column(db.String)
    areadesc = db.Column(db.Text)
    sent = db.Column(db.DateTime)
    effective = db.Column(db.DateTime)
    onset = db.Column(db.DateTime)
    expires = db.Column(db.DateTime)
    ends = db.Column(db.DateTime)
    status = db.Column(db.String)
    messagetype = db.Column(db.String)
    category = db.Column(db.String)
    severity = db.Column(db.String)
    certainty = db.Column(db.String)
    urgency = db.Column(db.String)
    event = db.Column(db.String)
    sender = db.Column(db.String)
    sendername = db.Column(db.String)
    headline = db.Column(db.Text)
    description = db.Column(db.Text)
    instruction = db.Column(db.Text)
    response = db.Column(db.String)
    
    def to_dict(self):
        """Convert model instance to dictionary."""
        return {
            'id': self.id,
            'alertid': self.alertid,
            'wxtype': self.wxtype,
            'areadesc': self.areadesc,
            'sent': self.sent.isoformat() if self.sent else None,
            'effective': self.effective.isoformat() if self.effective else None,
            'onset': self.onset.isoformat() if self.onset else None,
            'expires': self.expires.isoformat() if self.expires else None,
            'ends': self.ends.isoformat() if self.ends else None,
            'status': self.status,
            'messagetype': self.messagetype,
            'category': self.category,
            'severity': self.severity,
            'certainty': self.certainty,
            'urgency': self.urgency,
            'event': self.event,
            'sender': self.sender,
            'sendername': self.sendername,
            'headline': self.headline,
            'description': self.description,
            'instruction': self.instruction,
            'response': self.response
        }


    def __repr__(self):
        return f"<Alert {self.id} - {self.event} - {self.effective}>"

class Beaufort:
    def __init__(self, svg, text, css):
        self.svg = svg
        self.text = text
        self.css = css

class Astro:
    def __init__(self):
        self.sunrise: datetime = datetime.now()
        self.sunset: datetime = datetime.now()
        self.sunrise_tomorrow: datetime = datetime.now()
        self.sunset_tomorrow: datetime = datetime.now()
        self.darkness: timedelta = timedelta()
        self.daylight: timedelta = timedelta()
        self.elevation: float = 0.0
        self.has_sunset: bool = False

class Trend:
    """Class to represent a trend in weather measurements."""
    def __init__(self):
        self.trend = ""  # Direction of trend (up/down/Steady/Rising/Falling)
        self.by = 0.0    # Amount of change