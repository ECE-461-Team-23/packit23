# syntax=docker/dockerfile:1.4
FROM python:3.11-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY requirements.txt /app
RUN pip3 install -r requirements.txt

COPY . /app

ENV PORT 8080

# commenting because makes ci run indefinitely
ENTRYPOINT ["uvicorn"]
CMD ["src.main:app", "--host", "0.0.0.0", "--port", "8080"]