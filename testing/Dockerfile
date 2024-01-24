FROM python:3.8

WORKDIR /test

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

CMD ["python", "./testclient.py"]
