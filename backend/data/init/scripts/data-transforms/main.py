from sqlalchemy import create_engine
from merge_freq_and_word_info import get_full_word_info

CONNECTION_STRING = "postgresql+psycopg2://postgres:qwertyuiop@database-asmir.cmhmoaojrw66.ap-northeast-2.rds.amazonaws.com:5432/appdata"

def populate_table():
    words_pd = get_full_word_info()
    engine = create_engine(CONNECTION_STRING)
    words_pd.to_sql("words", engine)

populate_table()
