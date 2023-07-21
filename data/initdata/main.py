import sqlalchemy
import pandas as pd
from sqlalchemy import create_engine
from sqlalchemy.dialects import postgresql

import os
from dotenv import load_dotenv

from datatransforms.merge_freq_and_word_info import get_full_word_info
from datatransforms.simplify_definitions import simplify_definitions_column

load_dotenv()
CONNECTION_STRING = os.getenv("DB_STRING")

def is_words_table_empty(engine):
    countDf = pd.read_sql("""
                    SELECT COUNT(*) FROM korean_words;
                """, engine)
    return countDf["count"][0] == 0

def populate_table():
    engine = create_engine(CONNECTION_STRING)
    words_pd = get_full_word_info()

    words_pd.sort_values("freq_rank", inplace=True)

    words_pd["defs"] = words_pd["defs"].apply(lambda defs_list: tuple(defs_list))
    words_pd.drop_duplicates(subset=["word", "defs"], inplace=True)
    words_pd.reset_index(inplace=True, drop=True)

    print(words_pd)

    # simplify_definitions_column(words_pd)
    engine = create_engine(CONNECTION_STRING)
    words_pd.to_sql(
        "korean_words", 
        engine, 
        if_exists = "replace", # append just duplicates the words if they exists
        dtype = {
            "word": sqlalchemy.TEXT(), 
            "part_of_speech": sqlalchemy.TEXT(), 
            "defs": postgresql.ARRAY(sqlalchemy.TEXT()), 
            "freq_rank": sqlalchemy.INTEGER()
        }
    )

populate_table()
