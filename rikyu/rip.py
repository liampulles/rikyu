from pathlib import Path

from rikyu import error, step, docker
from rikyu.source import Source


class Ripper:
    def __init__(self, dest_base_path: Path, dvd_path: Path):
        self.dest_base_path = dest_base_path
        self.dvd_path = dvd_path
        self.titles = []

    def title(self, title_no: int, dest: Source):
        self.titles += [{"title_no": title_no, "dest": dest.relative_path()}]
        return self

    def rip(self):
        _verify_dir_exists(self.dest_base_path)
        _verify_dir_exists(self.dvd_path)
        step_defs = []
        for title in self.titles:
            dest = title["dest"]
            dest_abs = Path(self.dest_base_path, dest)
            title_no = title["title_no"]
            step_def = docker.docker_step_def(
                "lpulles/dvdtools:latest",
                [("dvd", self.dvd_path)],
                dest_abs,
                self.dest_base_path,
                args=["/sh/extractTitle.sh", str(title_no), "/in/dvd", "/out"]
            )
            step_defs += [step_def]
        step.run_steps(step_defs)


def _verify_dir_exists(path: Path):
    if not path.exists():
        raise error.RipperIOError(message=str(path) + " does not exist")
    if not path.is_dir():
        raise error.RipperIOError(message=str(path) + " is not a directory")


def _step_exec(path: Path):
    def _func():
        with Path(path, "extracted_test").open("w") as f:
            f.write("data!")

    return _func
