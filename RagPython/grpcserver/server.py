import grpc
from concurrent import futures
import pdfread_pb2
import pdfread_pb2_grpc

import sys
import os
# 获取当前文件的目录
current_dir = os.path.dirname(os.path.abspath(__file__))
# 获取 RagPython 目录的父目录
parent_dir = os.path.dirname(os.path.dirname(current_dir))
# 将父目录添加到 Python 路径
sys.path.append(parent_dir)

from RagPython.pdfreader.pdfread import TextVectorizer, open_and_read_pdf

class PdfServicer(pdfread_pb2_grpc.PdfServiceServicer):
    def __init__(self):
        self.vectorizer = TextVectorizer()

    def ExtractText(self, request, context):
        file_path = request.file_path
        try:
            pages_and_texts = open_and_read_pdf(file_path)
            for page in pages_and_texts:
                response = pdfread_pb2.PdfResponse(
                    page_number = page["page_number"],
                    sentence_chunk = page["sentence_chunk"],
                    chunk_char_count = page["chunk_char_count"],
                    chunk_word_count = page["chunk_word_count"],
                    chunk_token_count = page["chunk_token_count"]
                )
                response.embedding.extend(page["embedding"])  # 使用 extend 方法
                yield response

        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"Error extracting text: {str(e)}")

    def VectorizeText(self, request, context):
        try:
            vector = self.vectorizer.vectorize_text(request.text)
            return pdfread_pb2.VectorizeResponse(vector=vector)
        except ValueError as e:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details(str(e))
            return pdfread_pb2.VectorizeResponse()
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"向量化处理失败: {str(e)}")
            return pdfread_pb2.VectorizeResponse()





def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
   
    pdfread_pb2_grpc.add_PdfServiceServicer_to_server(
        PdfServicer(), 
        server)


    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()