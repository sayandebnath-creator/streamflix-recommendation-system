from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer

app = FastAPI()

model = SentenceTransformer("BAAI/bge-small-en-v1.5")


class EmbeddingRequest(BaseModel):
    text: str


@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/embed")
def embed(request: EmbeddingRequest):
    embedding = model.encode(request.text, normalize_embeddings=True)

    return {
        "embedding": embedding.tolist()
    }