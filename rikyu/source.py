from abc import abstractmethod
from pathlib import Path


class Source:
    @abstractmethod
    def relative_path(self) -> Path:
        pass


class UniqueSource(Source):
    def __init__(self, name: str):
        self.name = name

    def relative_path(self) -> Path:
        return Path(self.name)


class SequenceSource(Source):
    def __init__(self, name: str, idx: int):
        self.name = name
        self.idx = idx

    def relative_path(self) -> Path:
        return Path(self.name, str(self.idx))
