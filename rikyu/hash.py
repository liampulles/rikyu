import os
from hashlib import md5
from pathlib import Path
from typing import List, BinaryIO

from rikyu import error


def dir_hash(path: Path) -> str:
    sha = md5()
    files = _all_files_in(path)
    sorted(files, key=lambda f: f.as_posix())

    for f in files:
        sha.update(file_hash(path, f).encode())

    return sha.hexdigest()


def file_hash(base: Path, path: Path) -> str:
    try:
        megabyte = 1024 * 1024
        tot_bytes = path.stat().st_size
        sha = md5()
        rel = path.relative_to(base)
        sha.update(rel.as_posix().encode())
        with path.open("rb") as f:
            # (up to) first two megabytes
            sha.update(_hash_for_bytes(f, 2 * megabyte))
            if tot_bytes > 2 * megabyte:
                # (up to) last two megabytes. may be overlap, but that's fine
                f.seek(-2 * megabyte, os.SEEK_END)
                sha.update(_hash_for_bytes(f, 2 * megabyte))
        return sha.hexdigest()
    except Exception as e:
        raise error.HashIOError(e)


def _all_files_in(path: Path) -> List[Path]:
    return list(filter(lambda p: p.is_file(), path.glob("**/*")))


def _hash_for_bytes(f: BinaryIO, count: int) -> bytes:
    bytes_read = 0
    sha = md5()
    while bytes_read < count:
        count += 4096
        buf = f.read(4096)
        if not buf:
            break
        sha.update(md5(buf).digest())
    return sha.digest()
