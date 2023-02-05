from bs4 import BeautifulSoup
import urllib3, shutil


# Passing the source code to BeautifulSoup to create a BeautifulSoup object for it.
soup = BeautifulSoup(open("./Jaeger_UI.html"), "html.parser") 


# Extracting all the <a> tags into a list.
tags = soup.find_all('a', {'class': 'ResultItemTitle--item ub-flex-auto'})


def download_json(url):
    #json_url = 'http://localhost:16686/api/traces/' + url[-16::] + '?prettyPrint=true'
    json_url = 'http://10.42.1.8:16686/api/traces/' + url[-16::] + '?prettyPrint=true'
    print(json_url)
    http = urllib3.PoolManager()

    path = url[-16::] + '.json'

    with http.request('GET', json_url, preload_content=False) as r, open(path, 'wb') as out_file:
        shutil.copyfileobj(r, out_file)

# Extracting URLs from the attribute href in the <a> tags.
for tag in tags:
    link = tag.get('href')
    download_json(link)
