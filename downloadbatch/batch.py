from bs4 import BeautifulSoup
import urllib3, shutil


soup = BeautifulSoup(open("./Jaeger_UI.html"), "html.parser")

tags = soup.find_all('a', {'class': 'ResultItemTitle--item ub-flex-auto'})


def download_json(url):
    searcht = "trace/"
    tracei = url.rindex(searcht)    
    crop = tracei - len(url) + len(searcht)
    json_url = 'http://10.42.0.14:16686/api/traces/' + url[crop::] + '?prettyPrint=true'
    http = urllib3.PoolManager()

    path = url[crop::] + '.json'

    with http.request('GET', json_url, preload_content=False) as r, open(path, 'wb') as out_file:
        shutil.copyfileobj(r, out_file)

for tag in tags:
    link = tag.get('href')
    download_json(link)
