def extract_definitions(definitions):
    def extract_definition(definition):
        return definition["def"]
    return list(map(extract_definition, definitions))

def simplify_definitions_column(word_info_data_frame):
    word_info_data_frame["defs"] = list(map(extract_definitions, word_info_data_frame["defs"]))
