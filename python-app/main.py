from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from prometheus_fastapi_instrumentator import Instrumentator

app = FastAPI()

hello = {"en": "Hello", "fr": "Bonjour", "es": "Hola", "ml": "ഹലോ"}

# expose the default Python metrics to the /metrics endpoint
Instrumentator().instrument(app).expose(app)

@app.get("/{lang}")
async def get_hello(lang):
    return {"message": hello[lang]}
