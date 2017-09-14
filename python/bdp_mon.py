#encoding=utf8
import urllib.request
weburl = "http://123.59.41.67:57001/server/serverstate"  
webheader1 = {'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36'}  
webheader2 = {  
    'Connection': 'Keep-Alive',  
    'Accept': 'image/webp,image/apng,image/*,*/*;q=0.8',  
    'Accept-Language':'en-US,en;q=0.8', 
    "Cookie":"bdp_admin=eyJlbnYiOiJvbmxpbmUiLCJ1c2VybmFtZSI6ImJkcCIsImRpc3BsYXlOYW1lIjoiQkRQ55uR5o6n5bmz5Y+wIn0=; bdp_admin.sig=PVqPDflMWKIOYVR-QSynLeDIo1U",
    'User-Agent':'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.78 Safari/537.36',  
    'Accept-Encoding': 'gzip, deflate',  
    'Host': 'http://123.59.41.67:57001',  
    'DNT': '1'  
    }  
req = urllib.request.Request(url=weburl, headers=webheader2)    
webPage=urllib.request.urlopen(req)  
data = webPage.read()  
data = data.decode('utf-8')
print(data)  
print(type(webPage))  
print(webPage.geturl())  
print(webPage.info())  
print(webPage.getcode())
