import os
import json
import uuid
import tos

# 导入各模型 client
from audio.model.byte_asr import ByteASRClient

# 可扩展的模型 client 映射
MODEL_CLIENT_MAP = {
    "BYTEDANCE": ByteASRClient,
    # TODO: 添加其它模型 client
}

def load_config(config_path="../config.json"):
    with open(config_path, "r", encoding="utf-8") as f:
        return json.load(f)

def upload_to_tos(local_path, bucket, object_name, region, access_key, secret_key, endpoint):
    try:
        client = tos.TosClientV2(access_key, secret_key, endpoint, region)
        result = client.put_object_from_file(bucket, object_name, local_path)
        signed_url_obj = client.pre_signed_url(
            http_method=tos.HttpMethodType.Http_Method_Get,
            bucket=bucket,
            key=object_name,
            expires=3600
        )
        return signed_url_obj.signed_url
    except Exception as e:
        print(f'fail with error: {e}')
        raise

def is_local_file(path):
    return not (path.startswith("http://") or path.startswith("https://"))

def write_output(utterances, output_file):
    with open(output_file, "w", encoding="utf-8") as f:
        for seg in utterances:
            speaker = seg.get("additions", {}).get("speaker", "Unknown")
            text = seg.get("text", "").strip().replace("\n", "")
            if text:
                f.write(f'{speaker}: "{text}"\n')

def main():
    config = load_config(os.path.join(os.path.dirname(__file__), "../config.json"))
    model_src = config["AUDIO_CONFIGS"].get("MODEL_SRC", "BYTEDANCE").upper()
    if model_src not in MODEL_CLIENT_MAP:
        raise Exception(f"Unsupported model source: {model_src}")

    audio_configs = config["AUDIO_CONFIGS"][model_src]
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

    # 上传本地文件到TOS
    if is_local_file(audio_file):
        if not all([tos_bucket, tos_region, tos_access_key, tos_secret_key, tos_endpoint]):
            raise Exception("Tos Configuration is incomplete. Please check your config.json")
        object_name = f"audio2md/{uuid.uuid4().hex}_{os.path.basename(audio_file)}"
        print(f"[main] Uploading local file to TOS: {audio_file} -> {object_name}")
        audio_url = upload_to_tos(audio_file, tos_bucket, object_name, tos_region, tos_access_key, tos_secret_key, tos_endpoint)
        print(f"[main] Uploaded to: {audio_url}")
    else:
        audio_url = audio_file

    # 选择模型客户端
    model_client = MODEL_CLIENT_MAP[model_src]()

    print(f"[main] Submitting task to {model_src} API...")
    request_id, headers = model_client.submit_task(
        audio_url,
        app_key=app_key,
        access_key=access_key,
        uid=uid,
        audio_format=audio_format
    )
    print(f"[main] Task submitted, request_id: {request_id}")

    print("[main] Querying result...")
    result_json = model_client.query_task(request_id, headers)
    utterances = result_json.get("result", {}).get("utterances", [])
    if not utterances:
        print("[main] No utterances found in result.")
        return

    write_output(utterances, output_file)
    print(f"[main] Conversation saved to {output_file}")

if __name__ == "__main__":
    main()