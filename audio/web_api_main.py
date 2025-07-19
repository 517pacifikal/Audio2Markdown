import os
import json
import time
import uuid
import requests
import tos

def load_config(config_path="../config.json"):
    with open(config_path, "r", encoding="utf-8") as f:
        return json.load(f)

def upload_to_tos(local_path, bucket, object_name, region, access_key, secret_key, endpoint):
    try:
        client = tos.TosClientV2(access_key, secret_key, endpoint, region)
        result = client.put_object_from_file(bucket, object_name, local_path)
        print(f'http status code: {result.status_code}')
        print(f'request_id: {result.request_id}')
        print(f'crc64: {result.hash_crc64_ecma}')
        # 生成临时签名URL（有效期1小时），注意http_method要用枚举
        signed_url_obj = client.pre_signed_url(
        http_method=tos.HttpMethodType.Http_Method_Get,  # 用枚举类型
        bucket=bucket,
        key=object_name,
        expires=3600
    )
        signed_url = signed_url_obj.signed_url  # 取实际的URL字符串
        print(f'[upload_to_tos] Signed URL: {signed_url}')
        return signed_url
    except tos.exceptions.TosClientError as e:
        print(f'fail with client error, message:{e.message}, cause: {e.cause}')
        raise
    except tos.exceptions.TosServerError as e:
        print(f'fail with server error, code: {e.code}')
        print(f'error with request id: {e.request_id}')
        print(f'error with message: {e.message}')
        print(f'error with http code: {e.status_code}')
        print(f'error with ec: {e.ec}')
        print(f'error with request url: {e.request_url}')
        raise
    except Exception as e:
        print(f'fail with unknown error: {e}')
        raise

def submit_task(audio_url, app_key, access_key, uid, audio_format="mp3"):
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

def query_task(request_id, headers, max_wait=300):
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

def write_output(utterances, output_file):
    with open(output_file, "w", encoding="utf-8") as f:
        for seg in utterances:
            speaker = seg.get("additions", {}).get("speaker", "Unknown")
            text = seg.get("text", "").strip().replace("\n", "")
            if text:
                f.write(f'{speaker}: "{text}"\n')

def is_local_file(path):
    return not (path.startswith("http://") or path.startswith("https://"))

def main():
    config = load_config(os.path.join(os.path.dirname(__file__), "../config.json"))
    audio_configs = config["AUDIO_CONFIGS"]["OPENAI_API"]
    audio_file = audio_configs["AUDIO_FILE"]
    output_file = audio_configs["OUTPUT_FILE"]
    app_key = audio_configs["APP_KEY"]
    access_key = audio_configs["ACCESS_KEY"]
    uid = audio_configs.get("UID", "test_uid")
    audio_format = audio_configs.get("AUDIO_FORMAT", "mp3")


    # TOS配置
    tos_bucket = audio_configs.get("TOS_BUCKET")
    tos_region = audio_configs.get("TOS_REGION")
    tos_access_key = audio_configs.get("TOS_ACCESS_KEY")
    tos_secret_key = audio_configs.get("TOS_SECRET_KEY")
    tos_endpoint = audio_configs.get("TOS_ENDPOINT")

    # 如果是本地文件，先上传到TOS
    if is_local_file(audio_file):
        if not all([tos_bucket, tos_region, tos_access_key, tos_secret_key, tos_endpoint]):
            raise Exception("TOS配置缺失，无法上传本地音频文件！")
        object_name = f"audio2md/{uuid.uuid4().hex}_{os.path.basename(audio_file)}"
        print(f"[main] Uploading local file to TOS: {audio_file} -> {object_name}")
        audio_url = upload_to_tos(audio_file, tos_bucket, object_name, tos_region, tos_access_key, tos_secret_key, tos_endpoint)
        print(f"[main] Uploaded to: {audio_url}")
    else:
        audio_url = audio_file

    print("[main] Submitting task to API...")
    request_id, headers = submit_task(audio_url, app_key, access_key, uid, audio_format)
    print(f"[main] Task submitted, request_id: {request_id}")

    print("[main] Querying result...")
    result_json = query_task(request_id, headers)
    print(f"[main] Result: {json.dumps(result_json, ensure_ascii=False, indent=2)}")
    
    utterances = result_json.get("result", {}).get("utterances", [])
    if not utterances:
        print("[main] No utterances found in result.")
        return
    
    write_output(utterances, output_file)
    print(f"[main] 对话内容已输出到 {output_file}")

if __name__ == "__main__":
    main()