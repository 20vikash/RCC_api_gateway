import requests
import sys

def debug_code(code, language):
    prompt = f"Fix the bug in this code in {language} programming language if there is any and give only the code in triple backticks at the start without explanation\n{code}\n"

    OLLAMA_API_URL = "http://localhost:11434/api/generate"

    payload = {
        "model": "codellama:13b",
        "prompt": prompt,
        "stream": False
    }

    response = requests.post(OLLAMA_API_URL, json=payload)
    
    if response.status_code == 200:
        return response.json().get("response", "").strip()
    else:
        return f"Error: {response.status_code}, {response.text}"

args = sys.argv[1:]

code = debug_code(args[0], args[1])
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
