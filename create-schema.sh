#!/usr/bin/env bash
#
# Purpose: To create schema for stale-purger-db
#

set -euo pipefail

echo -n "Waiting for PostgreSQL to be ready "
while ! nc -z pg 5435; do
    echo -n "."
    sleep 0.1
done
echo " ok"

echo "Locating Connect backup file ... "
SCHEMA_FILE=${SCHEMA_FILE:-schema.sql}
if [[ -f $SCHEMA_FILE ]];
then
echo "schema.sql file located ... "
fi

echo "Creatign stale-purger database ... "
psql -hstale-purger -dstale-pods-info -Upostgres -p5435 -f $SCHEMA_FILE
echo "ok"