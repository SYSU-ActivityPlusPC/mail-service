FROM ubuntu

ADD main /
ENTRYPOINT [ "/main" ]