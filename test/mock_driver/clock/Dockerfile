FROM python:3.9-slim-buster
LABEL authors="AEnjoy"

WORKDIR /app
COPY . .
RUN pip install flask gunicorn gevent greenlet
EXPOSE 80

ENV FLASK_APP=clock.py
ENV FLASK_ENV=production

CMD ["gunicorn", "-b", ":80", "-w", "4", "--worker-class", "gevent", "clock:app"]
