FROM python:3.10-alpine

WORKDIR /microservice

COPY *.py /microservice/

COPY requirements.txt requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 50051

# (Optional) prevent Python from buffering output in Docker logs
ENV PYTHONUNBUFFERED=1

CMD [ "python", "app.py"]
