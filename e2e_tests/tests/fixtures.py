from datetime import datetime, timezone
import json
from urllib.parse import urlparse


ACCT_NAME = "acct:blog@localhost:3000"

PRIV_KEY_PEM = """-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQCni8P4bvkC/3Sx
NTrDw1qw0vWtJKZMsyJ3Mcs4+1apoVqOQhujUqGqFSiRT7Vmc7OEhB0vikdiTkCk
1PcoTd/rOa/0WbG8385JcLzaJfTIG+rrRtHwZ1TwxwBju43jpGEZjpbA0dqoXMkr
J1MyD7aPLoAiVe0ikw2czSZumv4ncemOtk0VG3b2TnIxo3CMKtUOWu8xT08MMIuo
3cZRnpI6Xr/ULzvjv8e3EjIpwRJqMPECtGsfdcjFmR0yFIrjrlmkZTiW31z/Dk7i
xRGD0ADy3WEQ3lA4l3mNZeyG4S0Wi4iYe9/wegESMZcakLoME7ks+KNS388Mdbcd
DKy9NmWvAgMBAAECggEABLQAA0hHhdWv6+Lc9xkpFuTvxTV4fuyvCf4u1eGlnstg
ZF/nW1/6w8XQ8WCgbJ4mKuZz1J14FYKxfoRaj8S9MA2Ff+wd+M77gRpAuDWajRzO
LQk8OW2yd7POXKkAzvln9F9eofkCFKR4zSpPGTenCJaQkuYrQEOKfUf7oofdRzQi
w9kmp3wAxM/EseHZpknYDCgDQV7MDQAaMD7kbynL2WfXPxebktwpRlKUwgtGrevj
gagQL8J/GX6wO3ymw9sln4BhlI2+3LuiMXQdQc1tamkXFCguCuOZCu/2VRdCHmiS
nnpu+FMspBHbvxO+RXo3Cu/S6jjJgoQxD2WZTE0gqQKBgQDM6AQdqBYjISdkI9Gl
6ZLLjwZRJSYpopujtX7pun61l9kUwQevaR2Z39rMWxX62DD6arazi/ygIUBw6Kgp
s/qBEb29ec+0cESdC8aJYb3dGvDzh/8C05p7ozxj8JZQcxq5W5jql/BELlSsUONO
jfqQv8RGZNSkD9uy6TxOr4eWIwKBgQDRUuO/XRDLt8Mp10mTshxTznSQ3gAJYKeG
0WfEC3kPEukHBQb8huqFcQDiQ71oBWuEdOQWgT3aBS6L+nIMyZMT5u+BejQm7/E5
pMM+z0VRpfFSsIrCvU8yKam0aemQGlKQAfhTct1gCg+wKnYsSQMlNHKWEfDbw9I/
cns/IN+dBQKBgQC6/Of0oFVDTZgC3GUPAO3C8QwUtM/0or1hUdk1Nck3shCZzeVT
f5tRtmSWpHCUbwGTJBsCEjdBcda6srXzCJkLe8Moy6ZtxR34KqzM5fM7eMB1nJ9s
Vunc9gPAN+cUF1ZF3H7ZZjoOHjGK5m3oW8xSl41np9Acv5P/2rP8Ilaa/QKBgQDJ
YwISfitGk8mEW8hB/L4cMykapztJyl/i6Vz31EHoKr1fL4sFMZg4QfwjtCBqD6zd
hshajoU/WHTr30wS2WxTXX9YBoZeX8KpPsdJioiagRioAYm+yfuDu2m2VZ+MMIb2
Xa7YOk6Zs5RcXL3M5YHNLaSAlUoxZTjGKhJBLhN1MQKBgQCbo3ngBl7Qjjx4WJ93
2WEEKvSDCv69eecNQDuKWKEiFqBN23LheNrN8DXMWFTtE4miY106dzQ0dUMh418x
K98rXSX3VvY4w48AznvPMKVLqesFjcvwnBdvk/NqXod20CMSpOEVj6W/nGoTBQt2
0PuW3IUym9KvO0WX9E+1Qw8mbw==
-----END PRIVATE KEY-----"""


def ensure_follow(client, inbox_url, actor_url):
    resp = client.post(
        inbox_url,
        json={
            "@context": "https://www.w3.org/ns/activitystreams",
            "id": "https://mock_masto/d0b5768b-a15b-4ed6-bc84-84c7e2b57588",
            "type": "Follow",
            "actor": "http://mock_masto:8000/users/h4kor",
            "object": actor_url,
        },
        headers={"Content-Type": "application/activity+json"},
    )
    assert resp.status_code == 200


def get_gmt_now() -> str:
    return datetime.now(datetime.UTC).strftime("%a, %d %b %Y %H:%M:%S GMT")


from http_message_signatures import (
    HTTPMessageSigner,
    HTTPMessageVerifier,
    HTTPSignatureKeyResolver,
    algorithms,
)
import requests, base64, hashlib, http_sfv


class MyHTTPSignatureKeyResolver(HTTPSignatureKeyResolver):
    keys = {"my-key": b"top-secret-key"}

    def resolve_public_key(self, key_id: str):
        return self.keys[key_id]

    def resolve_private_key(self, key_id: str):
        return priv_key
        # return PRIV_KEY_PEM


def sign(method, url, data):
    from urllib.parse import urlparse
    from cryptography.hazmat.primitives.serialization import load_pem_private_key
    from cryptography.hazmat.primitives import hashes
    from cryptography.hazmat.primitives.asymmetric import padding

    priv_key = load_pem_private_key(PRIV_KEY_PEM.encode(), None)
    body = json.dumps(data).encode()
    body_hash = hashlib.sha256(body).digest()
    digest = "SHA-256=" + base64.b64encode(body_hash).decode()
    date = datetime.now(tz=timezone.utc).strftime("%a, %d %b %Y %H:%M:%S GMT")
    host = "localhost:3000"
    target = urlparse(url).path
    to_sign = f"""(request-target): {method.lower()} {target}
host: {host}
date: {date}""".encode()
    sig = priv_key.sign(
        to_sign,
        padding.PKCS1v15(),
        hashes.SHA256(),
    )
    sig_str = base64.b64encode(sig).decode()

    request = requests.Request(method, url, data=body)
    request = request.prepare()
    request.headers["Content-Digest"] = digest
    request.headers["Host"] = host
    request.headers["Date"] = date
    request.headers["Signature"] = (
        f'keyId="http://mock_masto/users/h4kor#main-key",headers="(request-target) host date",signature="{sig_str}"'
    )
    return request
