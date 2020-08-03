FROM python:3.8-slim-buster

WORKDIR /dadid/sports-data-api/scraper-python
COPY ["master_dict.py", "scraper.py", "requirements.txt", "./"]
ENV PATH=/usr/lib/postgresql/X.Y/bin/:$PATH
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl && \
    curl https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb -o /chrome.deb && \
    dpkg -i /chrome.deb || apt-get install -yf && \
    rm /chrome.deb && \
    curl https://chromedriver.storage.googleapis.com/84.0.4147.30/chromedriver_linux64.zip -o ./chromedriver_linux64.zip && \
    apt-get install unzip && \
    unzip chromedriver_linux64.zip && \
    rm chromedriver_linux64.zip && \
    mv chromedriver /usr/local/bin/chromedriver && \
    chmod +x /usr/local/bin/chromedriver && \
    apt-get install -y libpq-dev && \
    apt-get install -y netcat-openbsd gcc && \
    apt-get clean && \
    python3 -m pip install -r requirements.txt --no-cache-dir
CMD ["python", "./scraper.py"]
