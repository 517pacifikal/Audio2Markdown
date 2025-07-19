import requests
import time
import uuid
from audio.model.model import ModelClient

class ByteASRClient(ModelClient):
    def submit_task(self, audio_url, app_key, access_key, uid, audio_format="mp3"):
        submit_url = "https://openspeech.bytedance.com/api/v3/auc/bigmodel/submit"
        request_id = str(uuid.uuid4())
        headers = {
            "Content-Type": "application/json",
            "X-Api-App-Key": app_key,
            "X-Api-Access-Key": access_key,
            "X-Api-Resource-Id": "volc.bigasr.auc",
            "X-Api-Request-Id": request_id,
            "X-Api-Sequence": "-1"
        }
        payload = {
            "user": {"uid": uid},
            "audio": {
                "format": audio_format,
                "url": audio_url,
            },
            "request": {
                "model_name": "bigmodel",
                "enable_itn": True,
                "enable_punc": True,
                "enable_speaker_info": True,
                "show_utterances": True
            }
        }
        resp = requests.post(submit_url, headers=headers, json=payload)
        if resp.headers.get("X-Api-Status-Code") != "20000000":
            raise Exception(f"Submit failed: {resp.headers.get('X-Api-Message')}")
        return request_id, headers

    def query_task(self, request_id, headers, max_wait=300):
        query_url = "https://openspeech.bytedance.com/api/v3/auc/bigmodel/query"
        headers = headers.copy()
        headers["X-Api-Request-Id"] = request_id
        for _ in range(max_wait):
            resp = requests.post(query_url, headers=headers, json={})
            status_code = resp.headers.get("X-Api-Status-Code")
            if status_code == "20000000":
                return resp.json()
            elif status_code in ("20000001", "20000002"):
                time.sleep(2)
                continue
            else:
                raise Exception(f"Query failed: {resp.headers.get('X-Api-Message')}")
        raise TimeoutError("Query timeout")