import pandas as pd
from dataclasses import make_dataclass

from datatransforms.extract_freq_info import extract_freq_info
from datatransforms.extract_word_info import extract_word_info

WordInfo = make_dataclass(
        "WordInfo", 
        [
            ("word", str), 
            ("part_of_speech", str), 
            ("defs", list[str])
        ])

FreqInfo = make_dataclass(
        "FreqInfo", 
        [
            ("word", str), 
            ("freq_rank", int)
        ])

def get_full_word_info():
    word_info_list = extract_word_info()
    freq_info_list = extract_freq_info()

    word_info_data_classes = map(lambda word_info: WordInfo(word_info["word"], word_info["part_of_speech"], word_info["defs"]), word_info_list)
    freq_info_data_classes = map(lambda freq_info: FreqInfo(freq_info["word"], freq_info["freq_rank"]), freq_info_list)

    word_info_frame = pd.DataFrame(word_info_data_classes)
    freq_info_frame = pd.DataFrame(freq_info_data_classes)

    full_info_frame = pd.merge(word_info_frame, freq_info_frame, on = "word")
    return full_info_frame