import mechanicalsoup
from tinydb import TinyDB, Query
from tinydb.operations import add, subtract
import json
browser = mechanicalsoup.StatefulBrowser()

def checkword(n):
    report = browser.open("https://pishugramotno.ru/morfologiya/%s" %(n)).raise_for_status()
    wordhttpdata = browser.page.find("h2", class_="morfologiya-word-block-title")
    if wordhttpdata != None:
    	return wordhttpdata.get_text()
    else:
    	return "Placeholder"

#frequency = {}
rustext = open("Gugo.txt", "r",encoding='utf-8')
txt = rustext.read().lower()
txt = txt.replace("\n", " ").replace("!", " ").replace(",", " ").replace(".", " ").replace("–", " ").replace("—", " ").replace(")", " ").replace("(", " ").replace("«", " ").replace("?", " ").replace("…", " ").replace("»", " ").split(" ")

rusray = list(filter(lambda a: a != '', txt))

# for word in rusray:
#     for double in frequency:
#         frequency.append({word, chckword(word), count})

db = TinyDB('db.json', encoding='utf-8', ensure_ascii=False)
SWord = db.table('RusWords')
for word in rusray:
    if SWord.get(Query().word == word) == None:
    	SWord.insert({'word': word, 'count': 1, 'stype': checkword(word)})
    else:
    	SWord.update(add('count', 1), Query().word == word)
    