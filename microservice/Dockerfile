FROM python:3.10-alpine

WORKDIR /microservice

COPY requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY . .

ENV FLASK_APP=app.py
ENV FLASK_RUN_PORT=8000
ENV FLASK_RUN_HOST=0.0.0.0

# (Optional) prevent Python from buffering output in Docker logs
ENV PYTHONUNBUFFERED=1

CMD [ "python3", "-m" , "flask", "run"]
