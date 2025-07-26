#!/usr/bin/env bash
#
# Purpose: Dump databases with the Data.
#

set -euo pipefail

#
# Wait until pg starts actually responding to its port; --wait doesn't appear to fully handle this situation
# and sometimes we end up in races.
#
echo -n "Waiting for PostgreSQL to be ready "
while ! nc -z pg 5435; do
    echo -n "."
    sleep 0.1
done
echo " ok"

echo "Locating backup file ... "
DUMP_FILE_NAME=${DUMP_FILE_NAME:-db_dump.dump}
if [[ -f $DUMP_FILE_NAME ]];
then
echo "Connect backup file located ... "
fi

echo "Resetting DB database ... "
pg_restore --clean --if-exists --no-owner -h -dstale-pods-info -Upostgres -p5435 "$DUMP_FILE_NAME"
echo "ok"

echo "Database reset complete."