import yaml

KOR_WORDS_PATH = "./../../original-data/korean-words-basic-info.yml"
PART_OF_SPEECH_PATH = "./../../original-data/part-of-speech.yml"

def extract_word_info(word_info_path = KOR_WORDS_PATH, part_of_speech_path = PART_OF_SPEECH_PATH):
    with open(word_info_path) as word_info, open(part_of_speech_path) as part_of_speech:
        word_info_text = word_info.read()
        part_of_speech_text = part_of_speech.read()

        list_of_words = yaml.safe_load(word_info_text)
        part_of_speech = yaml.safe_load(part_of_speech_text)
        
        list_of_words = map(lambda word: {"word": word["word"], "pos": word["pos"], "defs": word["defs"]}, list_of_words)
        list_of_words = map(lambda word: {"part_of_speech": part_of_speech[word["pos"]], "word": word["word"], "defs": word["defs"]}, list_of_words)

        return list(list_of_words)
