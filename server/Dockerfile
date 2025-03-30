FROM python:slim

RUN pip install -r requirements.txt
RUN pip install gunicorn

COPY . .
RUN chmod a+x boot.sh

EXPOSE 3000
ENTRYPOINT ["./boot.sh"]
