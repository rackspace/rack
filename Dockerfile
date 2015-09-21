FROM ubuntu:15.04

RUN apt-get update -y && apt-get install vim git -y

ADD . /rack

WORKDIR /rack

CMD bash -c "source script/lib.sh && update_docs 1.0.woo && git diff"
