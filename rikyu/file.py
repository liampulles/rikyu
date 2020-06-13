from pathlib import Path

import jsonpickle

from rikyu import error


def load_obj(path: Path):
    try:
        with path.open("r") as f:
            json = f.read()
    except Exception as e:
        raise error.ConfigIOError(e)
    return jsonpickle.decode(json)


def save_obj(obj, path: Path):
    json = jsonpickle.encode(obj)
    try:
        with path.open("w") as f:
            f.write(json)
    except Exception as e:
        raise error.ConfigIOError(e)
