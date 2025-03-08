import requests
import sys

def generate_code(prompt, language="Python"):
    OLLAMA_API_URL = "http://localhost:11434/api/generate"

    payload = {
        "model": "codellama:13b",
        "prompt": f"{prompt}.\nGenerate a code for this in {language}.",
        "stream": False
    }

    response = requests.post(OLLAMA_API_URL, json=payload)
    
    if response.status_code == 200:
        return response.json().get("response", "").strip()
    else:
        return f"Error: {response.status_code}, {response.text}"

args = sys.argv[1:]

code = generate_code(args[0], args[1])
lines = code.split("\n")

res = []
record = False

for line in lines:
    if line.startswith("```"):
        record = not record
        continue
    
    if record:
        res.append(line)

print("\n".join(res))
