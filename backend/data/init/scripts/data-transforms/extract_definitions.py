def extract_definitions(definitions):
    def extract_definition(definition):
        return definition["def"]
    return list(map(extract_definition, definitions))