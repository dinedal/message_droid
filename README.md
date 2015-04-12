message_droid
=============

API
---

### /update

Add/update your service's message.

```bash
curl -i -X POST -d'{"service_id": "your_unique_service_name", "text": "your message here"}' http://10.0.0.223:8080/update

HTTP/1.1 200 OK
Content-Length: 0
Content-Type: text/plain; charset=utf-8
```
