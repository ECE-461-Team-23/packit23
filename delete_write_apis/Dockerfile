# syntax=docker/dockerfile:1.4
FROM python:3.8-alpine

WORKDIR /app

COPY requirements.txt /app
RUN pip3 install -r requirements.txt

COPY . /app

ENV PORT 8080

RUN pip3 install -e .

# commenting because makes ci run indefinitely
# ENTRYPOINT ["flask"]
# CMD ["--app", "src", "run", "-p", "8080"]