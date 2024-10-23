import fitz  # PyMuPDF
from tqdm.auto import tqdm
import os

def text_formatter(text: str) -> str:
    """Performs minor formatting on text."""
    cleaned_text = text.replace("\n", " ").strip() # note: this might be different for each doc (best to experiment)
    # Other potential text formatting functions can go here
    return cleaned_text

# Open PDF and get lines/pages
# Note: this only focuses on text, rather than images/figures etc
def open_and_read_pdf(pdf_path: str, max_pages: int = 20) -> list[dict]:

    # 验证PDF文件是否存在
    if not os.path.exists(pdf_path):
        raise FileNotFoundError(f"PDF文件不存在: {pdf_path}")
    doc = fitz.open(pdf_path)  # open a document
    pages_and_texts = []
    for page_number, page in tqdm(enumerate(doc)):  # iterate the document pages
        if page_number >= max_pages:
            break
        text = page.get_text()  # get plain text encoded as UTF-8
        text = text_formatter(text)
        pages_and_texts.append({"page_number": page_number ,  # adjust page numbers since our PDF starts on page 42
                                "page_char_count": len(text),
                                "page_word_count": len(text.split(" ")),
                                "page_sentence_count_raw": len(text.split(". ")),
                                "page_token_count": len(text) // 4,  # 1 token = ~4 chars, see: https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them
                                "text": text})
        
    
    pages_and_texts = get_nlp_pdf_pages_and_texts(pages_and_texts)
   
    num_sentence_chunk_size = 10 
    for item in tqdm(pages_and_texts):
        item["sentence_chunks"] = split_list(input_list=item["sentences"], slice_size=num_sentence_chunk_size)
        item["num_chunks"] = len(item["sentence_chunks"])

    pages_and_chunks = get_and_filter_pdf_chunks(pages_and_texts)
    pages_and_chunks = get_pdf_chunks_with_embeddings(pages_and_chunks)
   
    return pages_and_chunks



import pandas as pd
def get_pdf_stats(pages_and_texts: list[dict]) -> tuple[pd.DataFrame, pd.DataFrame]:
   
    df = pd.DataFrame(pages_and_texts)
    return df.head(), df.describe().round(2)

def get_nlp_pdf_pages_and_texts(pages_and_texts: list[dict]) -> list[dict]:
    from spacy.lang.en import English 
    nlp = English()
    nlp.add_pipe("sentencizer")
    for item in tqdm(pages_and_texts):
        item["sentences"] = list(nlp(item["text"]).sents)
        # Make sure all sentences are strings
        item["sentences"] = [str(sentence) for sentence in item["sentences"]]
        # Count the sentences 
        item["page_sentence_count_spacy"] = len(item["sentences"])
    return pages_and_texts
# Define split size to turn groups of sentences into chunks

# Create a function that recursively splits a list into desired sizes
def split_list(input_list: list, 
               slice_size: int) -> list[list[str]]:
    """
    Splits the input_list into sublists of size slice_size (or as close as possible).

    For example, a list of 17 sentences would be split into two lists of [[10], [7]]
    """
    return [input_list[i:i + slice_size] for i in range(0, len(input_list), slice_size)]
# Loop through pages and texts and split sentences into chunks

def get_and_filter_pdf_chunks(pages_and_texts: list[dict]) -> list[dict]:
    import re
    # Split each chunk into its own item
    pages_and_chunks = []
    for item in tqdm(pages_and_texts):
        for sentence_chunk in item["sentence_chunks"]:
            chunk_dict = {}
            chunk_dict["page_number"] = item["page_number"]       
            # Join the sentences together into a paragraph-like structure, aka a chunk (so they are a single string)
            joined_sentence_chunk = "".join(sentence_chunk).replace("  ", " ").strip()
            joined_sentence_chunk = re.sub(r'\.([A-Z])', r'. \1', joined_sentence_chunk) # ".A" -> ". A" for any full-stop/capital letter combo 
            chunk_dict["sentence_chunk"] = joined_sentence_chunk

            # Get stats about the chunk
            chunk_dict["chunk_char_count"] = len(joined_sentence_chunk)
            chunk_dict["chunk_word_count"] = len([word for word in joined_sentence_chunk.split(" ")])
            chunk_dict["chunk_token_count"] = len(joined_sentence_chunk) // 4 # 1 token = ~4 characters
        
            pages_and_chunks.append(chunk_dict)
        
    df = pd.DataFrame(pages_and_chunks)
    #过滤小于30token的chunk
    min_token_length = 30
    pages_and_chunks_over_min_token_len = df[df["chunk_token_count"].astype(int)  > min_token_length].to_dict(orient="records")
    return pages_and_chunks_over_min_token_len

import logging
def get_embedding(text: str) -> dict[str, list[float]]:
    from sentence_transformers import SentenceTransformer
    embedding_model = SentenceTransformer(model_name_or_path="all-mpnet-base-v2", 
                                device="cpu")
    embeddings = embedding_model.encode(text,batch_size=int(32),convert_to_tensor=True)
    logging.debug(f"Original embeddings shape: {embeddings.shape}")
    logging.debug(f"Embeddings type: {type(embeddings)}")
    if embeddings.ndim == 1:
        embeddings = embeddings.reshape(1, -1)
    embeddings_list = embeddings.tolist()   
    if isinstance(text, str):
        text = [text]
    result = {t: [float(e) for e in emb] for t, emb in zip(text, embeddings_list)}
    logging.debug(f"Result keys: {list(result.keys())}")
    logging.debug(f"First embedding length: {len(next(iter(result.values())))}")
   
    return result

def get_pdf_chunks_with_embeddings(pages_and_chunks: list[dict]) -> list[dict]:
    for item in tqdm(pages_and_chunks):
        embedding_result = get_embedding(item["sentence_chunk"])
        if not embedding_result:
            logging.warning(f"Empty embedding result for chunk: {item['sentence_chunk'][:50]}...")
            continue
        item["embedding"] = next(iter(embedding_result.values()))  # 获取字典中的第一个（也是唯一的）值
    return pages_and_chunks


