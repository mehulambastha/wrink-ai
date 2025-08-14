import logging
import grpc
import aiohttp
import asyncio
import requests
from concurrent import futures
import search_pb2
import search_pb2_grpc
from bs4 import BeautifulSoup


logging.basicConfig(level=logging.INFO)

class Searcher(search_pb2_grpc.SearchServiceServicer):
    def simple_scrape(self, url: str) -> str:
        try:
            headers = {
                'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
            }
            response = requests.get(url, headers=headers, timeout=10)
            response.raise_for_status()

            soup = BeautifulSoup(response.text, 'html.parser')
            return soup.get_text()
        except requests.RequestException as e:
            logging.error(f"Error fetching {url}: {e}")
            return f"Error fetching content from {url}."


    def Search(self, request, context):
        logging.info(f"Recieved Search request for topic: {request.topic}, keywords: {request.keywords}")

        search_query = "+".join(request.keywords)
        search_url = f"https://www.google.com/search?q={search_query}"

        logging.info(f"Scariping URL: {search_url}")
        scraped_text = self.simple_scrape(search_url)

        processed_text = (scraped_text[:1000] + '...') if len(scraped_text) > 1000 else scraped_text

        logging.info("Search and processing complete. Sendingn response.")

        return search_pb2.SearchResponse(search_results_text=processed_text)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    search_pb2_grpc.add_SearchServiceServicer_to_server(Searcher(), server)
    server.add_insecure_port('[::]:50051')
    logging.info("Python Search gRPC server starting on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()

