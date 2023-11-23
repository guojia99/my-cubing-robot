import argparse
import time

from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service


# python load_mycube_image.py --image 2.png --url https://mycube.club/
# 使用该脚本, 你必须去下载两样东西
# chromedriver 和 chrome
# 1. 版本必须一致, 且大于114.x
# 2. 放在该脚本同等的位置上

def main():
    parser = argparse.ArgumentParser(description="image_load")
    parser.add_argument("--image", help="path")
    parser.add_argument("--url", help="url")
    parser.add_argument("--wait", help="wait time", default=5)

    args = parser.parse_args()
    if args.image and args.url:
        service = Service(executable_path="/usr/local/bin/chromedriver-linux64/chromedriver")
        chrome_options = Options()
        chrome_options.add_argument('--headless')
        chrome_options.add_argument('--disable-gpu')
        chrome_options.add_argument('--start-maximized')
        chrome_options.add_argument("--no-sandbox")
        chrome_options.binary_location = "/usr/local/bin/chrome-linux64/chrome"
        driver = webdriver.Chrome(service=service, options=chrome_options)
        driver.get(args.url)
        driver.maximize_window()

        time.sleep(args.wait)
        driver.fullscreen_window()

        height = driver.execute_script("return document.documentElement.scrollHeight")

        driver.set_window_size(1980, height / 4 * 3)
        driver.save_screenshot(args.image)
        driver.quit()


if __name__ == "__main__":
    main()
