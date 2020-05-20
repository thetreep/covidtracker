FROM golang:1.14


WORKDIR /app

# copy the linux executable
COPY covidtracker covidtracker
RUN chmod +x /app/covidtracker

# Env variables

EXPOSE 3456

ENTRYPOINT ["/app/covidtracker"]
