from bs4 import BeautifulSoup

PATH_TO_FREQ_LIST = "./../../data/original/wiki-5800.html"

def extract_freq_info(file_path = PATH_TO_FREQ_LIST):
    with open(file_path) as file:
        file_text = file.read()

        soup = BeautifulSoup(file_text, "html.parser")
        table = soup.table
        
        result = []
        for col in table.find_all("td"):
            freq_range = col.div.get_text()
            range_start = int(freq_range.split()[0])# format 'start_index - end_index'

            for word_index, word_link in enumerate(col.find_all("a")):
                word = word_link.get_text()
                freq_rank = range_start + word_index
                result.append({"word": word, "freq_rank": freq_rank})

        return result
