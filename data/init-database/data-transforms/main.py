#!./../../venv/bin/python3
import sqlalchemy
from sqlalchemy import create_engine
from sqlalchemy.dialects import postgresql

import os
from dotenv import load_dotenv

from merge_freq_and_word_info import get_full_word_info
from simplify_definitions import simplify_definitions_column

load_dotenv()
# CONNECTION_STRING = "postgresql+psycopg2://postgres:qwertyuiop@database-asmir.cmhmoaojrw66.ap-northeast-2.rds.amazonaws.com:5432/appdata"
CONNECTION_STRING = os.getenv("DB_STRING")


def populate_table():
    words_pd = get_full_word_info()
    simplify_definitions_column(words_pd)
    engine = create_engine(CONNECTION_STRING)
    words_pd.to_sql("korean_words", engine, if_exists = "replace", dtype = {"word": sqlalchemy.TEXT(), "part_of_speech": sqlalchemy.TEXT(), "defs": postgresql.ARRAY(sqlalchemy.TEXT()), "freq_rank": sqlalchemy.INTEGER()})

populate_table()
