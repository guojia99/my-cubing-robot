from selenium import webdriver
import time
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service

service = Service(executable_path="./chromedriver-linux64/chromedriver")
chrome_options = Options()
chrome_options.add_argument('--headless')
chrome_options.add_argument('--disable-gpu')
chrome_options.add_argument('--start-maximized')
chrome_options.binary_location = "./chrome-linux64/chrome"

driver = webdriver.Chrome(service=service, options=chrome_options)
driver.get("https://mycube.club/contest?id=36&contest_tab=tab_nav_all_score_table")
driver.maximize_window()

time.sleep(5)
driver.fullscreen_window()

height = driver.execute_script("return document.documentElement.scrollHeight")

driver.set_window_size(1980, height / 4 * 3)
driver.save_screenshot("1.png")
driver.quit()



