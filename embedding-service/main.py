from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer

app = FastAPI()

model = SentenceTransformer("BAAI/bge-small-en-v1.5")


class EmbeddingRequest(BaseModel):
    text: str
    is_query: bool = False


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/embed")
def embed(request: EmbeddingRequest):
    text = request.text

    if request.is_query:
        text = f"Represent this sentence for searching relevant passages: {text}"

    embedding = model.encode(
        text,
        normalize_embeddings=True,
    )

    return {
        "embedding": embedding.tolist()
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)