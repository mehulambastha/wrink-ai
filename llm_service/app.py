import logging
import grpc
import os
from concurrent import futures
import google.generativeai as ai

from dotenv import load_dotenv
load_dotenv()

import llm_pb2
import llm_pb2_grpc

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class LLMGenerator(llm_pb2_grpc.LLMServiceServicer):

    def __init__(self):
        self.api_key = os.getenv("GEMINI_API_KEY")
        if not self.api_key:
            raise ValueError("Missing API key")

        ai.__builtins__.

    
