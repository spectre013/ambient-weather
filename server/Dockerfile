FROM python:slim

COPY . .
RUN pip install -r requirements.txt
RUN pip install gunicorn

RUN chmod a+x boot.sh

EXPOSE 3000
ENTRYPOINT ["./boot.sh"]
