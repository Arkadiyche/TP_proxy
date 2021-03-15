FROM golang:1.13-stretch AS builder

WORKDIR /build
COPY . .

USER root
RUN go build  ./cmd/main.go

FROM ubuntu:20.04
COPY . .

EXPOSE 5432
EXPOSE 8080

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get -y update && apt -y install postgresql-12

USER postgres

RUN  /etc/init.d/postgresql start &&\
    psql --command "CREATE USER secure WITH SUPERUSER PASSWORD 'secure';" &&\
    createdb -O secure secure &&\
    psql -f postgre.sql -d secure &&\
    /etc/init.d/postgresql stop


RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/12/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/12/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]



#USER postgres
#CMD ["/usr/lib/postgresql/12/bin/postgres", "-D", "/var/lib/postgresql/12/main", "-c", "config_file=/etc/postgresql/12/main/postgresql.conf"]
#CMD ./run

USER root
COPY --from=builder  /build/main /usr/bin
CMD /etc/init.d/postgresql start && main