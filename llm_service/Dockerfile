FROM python:3.10-alpine

WORKDIR /app

COPY *.py /app/

COPY requirements.txt .
RUN pip install --no-chache-dir -r requirements.txt

COPY app.py .

EXPOSE 50052

CMD ["python", "app.py"]
