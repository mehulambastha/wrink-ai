FROM python:3.10-alpine

WORKDIR /microservice

COPY search_pb2.py .
COPY search_pb2_grpc.py .

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

EXPOSE 50051

# (Optional) prevent Python from buffering output in Docker logs
ENV PYTHONUNBUFFERED=1
ENV NEWSAPI=aed78dda66cf48de8464d1153ac30698

CMD ["watchmedo", "auto-restart", "--directory=.", "--pattern=*.py", "--recursive", "--", "python", "app.py"]
