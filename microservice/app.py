import logging
import time
import os
import grpc
from urllib.parse import quote_plus
from concurrent import futures
import search_pb2
import search_pb2_grpc
from bs4 import BeautifulSoup
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, WebDriverException

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class Searcher(search_pb2_grpc.SearchServiceServicer):
    def __init__(self):
        self.selenium_url = os.getenv("SELENIUM_HUB_URL")
        if not self.selenium_url:
            raise ValueError("SELENIUM_HUB_URL")
        logging.info(f"Selenium HUb running on the url {self.selenium_url}")


    def get_driver(self):
        chrome_options = webdriver.ChromeOptions()
        
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-blink-features=AutomationControlled")
        chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
        chrome_options.add_experimental_option('useAutomationExtension', False)
        chrome_options.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

        try:
            driver = webdriver.Remote(
                command_executor=self.selenium_url,
                options=chrome_options
            )
            return driver
        except WebDriverException as e:
            logging.error(f"Failed to create new Selenium session due to: {e}")
            return None


    def scrape_article_text(self, driver, url: str):
        try:
            logging.info(f"Scraping from timesofindia.indiatimes.com")
            driver.get(url)

            wait = WebDriverWait(driver, 20)

            # TODO: FIX THE SCRAPING FOR TOI, FIND OUT IF TOI RESTRICTS SCRAPING. FIND OUT WHAT DOES THE DRIVER SEE.
            article_body = wait.until(EC.presence_of_element_located((By.CLASS_NAME, "_s30J clearfix  ")))

            soup = BeautifulSoup(article_body.get_attribute('innerHTML'), 'html.parser')

            for script_or_style in soup(["script", "style"]):
                script_or_style.decompose()

            return soup.get_text(separator=' ', strip=True)

        except (TimeoutException, WebDriverException) as e:
            logging.error(f"Error scraping the article {url}: {e}")
            return f"Could not scrape content from {url}"

    def Search(self, request, context):
        """
        Handle Single gRPC request by creating a new selenium session for this request.    
        """
        logging.info(f"Processing saerch request for topic: {request.topic}")

        driver = self.get_driver()
        if not driver:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details("Could not create Selenium browser session")
            return search_pb2.SearchResponse()

        try:
            # Lets go to google news first
            query = quote_plus(" ".join(request.keywords) + " site:timesofindia.indiatimes.com")
            search_url = f"https://duckduckgo.com/?q={query}&t=h_&ia=web"

            driver.get("about:blank")
            driver.execute_script("Object.defineProperty(navigator, 'webdriver', {get: () => undefined})")

            logging.info((f"Navigating to: {search_url}"))
            driver.get(search_url)

            time.sleep(2)
            logging.info("Page title: " + driver.title)
                
            wait = WebDriverWait(driver, 20) # max wait time
            link_elements = wait.until(EC.presence_of_all_elements_located((By.CSS_SELECTOR, "a[data-testid='result-title-a']")))

            logging.info(f"Found {link_elements} link elements")
            
            article_urls = [el.get_attribute('href') for el in link_elements[]]

            # TODO: ITERETE OVER ARTICLES BY ONLY SELECTING TOI ONES, BASICALLY SKIPPPING THE ADVERTS.

            # for h3 in link_elements:
            #     try:
            #         parent_link = h3.find_element(By.XPATH, "./acestor::a")
            #         url=parent_link.get_attribute('href')
            #
            #         if url and url.startswith('http'):
            #             article_urls.append(url)
            #
            #         if len(article_urls) >= 5:
            #             break
            #
            #     except Exception:
            #         continue

            if not article_urls:
                return search_pb2.SearchResponse(search_results_text="No articles found")

            logging.info(f"Found {len(article_urls)} article links")

            all_content = []

            for url in article_urls:
                content = self.scrape_article_text(driver, url)
                all_content.append(content)

            final_text = "\n\n--- NEW ARTICLE --- \n\n".join(all_content)
            return search_pb2.SearchResponse(search_results_text=final_text)
        except (TimeoutException, WebDriverException) as e:
            logging.error(f"A Selenium error occuered: {e}")
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(f"A Selenium error occuered: {e}")
            return search_pb2.SearchResponse()
        finally:
            if driver:
                driver.quit()
                logging.info("Selenium session closed.")



def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    search_pb2_grpc.add_SearchServiceServicer_to_server(Searcher(), server)
    server.add_insecure_port('[::]:50051')
    logging.info("Python Search gRPC server starting on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()

