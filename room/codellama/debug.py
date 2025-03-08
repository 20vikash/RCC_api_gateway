import requests
import sys

def debug_code(code):
    prompt = f"Fix the bug in this code if there is any\n{code}\n"

    OLLAMA_API_URL = "http://localhost:11434/api/generate"

    payload = {
        "model": "codellama:latest",
        "prompt": prompt,
        "stream": False
    }

    response = requests.post(OLLAMA_API_URL, json=payload)
    
    if response.status_code == 200:
        return response.json().get("response", "").strip()
    else:
        return f"Error: {response.status_code}, {response.text}"

args = sys.argv[1:]

code = debug_code(args[0])
lines = code.split("\n")

res = ""

record = False
for line in lines:
    if line == "```" and not record:
        record = True
        continue
    elif line == "```" and record:
        record = False
        break
    
    if record:
        res += line + "\n"

print(res)
