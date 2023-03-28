# syntax=docker/dockerfile:1.4
FROM python:3.8-alpine

WORKDIR /app

COPY requirements.txt /app
RUN pip3 install -r requirements.txt

COPY . /app

ENV PORT 8080

ENTRYPOINT ["python3"]
CMD ["app.py"]
