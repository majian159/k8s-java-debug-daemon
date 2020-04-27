FROM ubuntu:xenial
WORKDIR /app
COPY jdd-linux ./jdd
COPY craw.sh ./
RUN chmod +x jdd && mkdir stacks
CMD ./jdd