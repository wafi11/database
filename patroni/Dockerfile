FROM postgres:14-bullseye

RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 \
    python3-pip \
    python3-setuptools \
    curl \
    vim \
    && pip3 install --no-cache-dir \
    patroni[etcd3]==3.3.0 \
    psycopg2-binary \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/patroni /data/patroni \
    && chown -R postgres:postgres /etc/patroni /data/patroni

EXPOSE 5432 8008

USER postgres

ENTRYPOINT ["patroni", "/etc/patroni/patroni.yml"]
