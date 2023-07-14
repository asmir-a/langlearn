set -a
. ./../backend/.env
set +a

psql $DB_STRING -f ./migrations/initial-tables.sql

python3 -m pip install -r ./initdata/requirements.txt
python3 ./initdata/main.py
