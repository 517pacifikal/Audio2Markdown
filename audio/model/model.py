from abc import ABC, abstractmethod


class ModelClient(ABC):
    @abstractmethod
    def submit_task(self, audio_url, **kwargs):
        pass

    @abstractmethod
    def query_task(self, request_id, headers, max_wait=300):
        pass