set -a
. ./../backend/.env
export DB_STRING=postgresql://postgres:qwertyuiop@langlearndb.cmhmoaojrw66.ap-northeast-2.rds.amazonaws.com/langlearn

set +a

psql $DB_STRING -f ./migrations/initial-tables.sql

python3 -m pip install -r ./initdata/requirements.txt
python3 ./initdata/main.py
