import os
import yaml

current_source_file = os.path.dirname(__file__)
KOR_WORDS_PATH = os.path.join(current_source_file, "./../originaldata/korean-words-basic-info.yml")
PART_OF_SPEECH_PATH = os.path.join(current_source_file, "./../originaldata/part-of-speech.yml")

def is_word_valid(word):
    if word["word"] == None or word["defs"] == None or word["pos"] == None:
        return False
    return True

def map_word_from_yaml_to_db(word):
    def map_word_def_from_yaml_to_db(word_def):
        return word_def["def"]
    return {
        "word": word["word"],
        "pos": word["pos"],
        "defs": list(map(map_word_def_from_yaml_to_db, word["defs"]))
    }

def extract_word_info(word_info_path = KOR_WORDS_PATH, part_of_speech_path = PART_OF_SPEECH_PATH):
    with open(word_info_path) as word_info, open(part_of_speech_path) as part_of_speech:
        word_info_text = word_info.read()
        part_of_speech_text = part_of_speech.read()

        list_of_words = yaml.safe_load(word_info_text)
        part_of_speech = yaml.safe_load(part_of_speech_text)
        
        list_of_words = filter(is_word_valid, list_of_words)
        list_of_words = map(map_word_from_yaml_to_db, list_of_words)
        list_of_words = map(lambda word: {"part_of_speech": part_of_speech[word["pos"]], "word": word["word"], "defs": word["defs"]}, list_of_words)
        return list(list_of_words)