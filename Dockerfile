FROM ubuntu:xenial
WORKDIR /app
COPY jdd ./
RUN chmod +x jdd
CMD ./jdd