import grpc
import logging

import search_pb2
import search_pb2_grpc

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def run_test():
    logging.info("Attemptintg to connect to grpc client")

    with grpc.insecure_channel('localhost:50051') as channel:
        stub = search_pb2_grpc.SearchServiceStub(channel)

        logging.info("--Sending seach request--")

        request = search_pb2.SearchRequest(
            topic="Tea App Scandal",
            keywords=["Tea App", "cybersecurity", "data breach"]
        )

        try:
            response = stub.Search(request, timeout=60)
            logging.info("--Recieved response--")
            with open('htmldump.txt', 'w') as f:
                f.write(str(response))
            print(f"File written")
        except grpc.RpcError as e:
            logging.error(f"Failed to receive response  due to: {e}")

if __name__ == "__main__":
    run_test()
