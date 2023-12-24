FROM python:3.9

WORKDIR /app

RUN pip install --upgrade pip

COPY ./python/requirements.txt requirements.txt

RUN pip install -r requirements.txt

COPY ./python/app.py app.py

CMD ["faust", "-A", "app", "worker", "-l", "info"]
