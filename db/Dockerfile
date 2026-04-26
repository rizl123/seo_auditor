FROM postgres:16-alpine

ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=postgres

COPY init.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init.sh

HEALTHCHECK --interval=5s --timeout=5s --retries=5 \
    CMD pg_isready -U postgres
