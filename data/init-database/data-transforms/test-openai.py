import openai

openai.api_key = "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"
result = openai.Model.list()

print(result)

response = openai.Image.create(
        prompt="a white dog",
        n = 1,
        size = "1024x1024"
)

print(response)
