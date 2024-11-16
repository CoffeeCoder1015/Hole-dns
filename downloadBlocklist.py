import requests

urls = [
    "https://gist.githubusercontent.com/anudeepND/adac7982307fec6ee23605e281a57f1a/raw/5b8582b906a9497624c3f3187a49ebc23a9cf2fb/Test.txt",
    "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
]

resp = requests.get(urls[1])
with open("default_host_raw.txt","w") as fio:
    fio.write(resp.text)