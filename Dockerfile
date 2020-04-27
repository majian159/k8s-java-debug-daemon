FROM ubuntu:xenial
WORKDIR /app
COPY jdd-linux ./jdd
RUN chmod +x jdd && mkdir stacks
CMD ./jdd