from pathlib import Path

from rikyu import error, file, rip

_project_file = ".rikyu.project"


class Project:
    def __new__(cls, path: Path):
        try:
            return file.load_obj(path.joinpath(_project_file))
        except error.ConfigIOError:
            self = object.__new__(cls)
            _make_folder(path)
            self.path = path
            file.save_obj(self, path.joinpath(_project_file))
            return self

    def rip_from_folder(self, dvd_path: Path) -> rip.Ripper:
        return rip.Ripper(self.path, dvd_path)


def _make_folder(path: Path):
    try:
        path.mkdir(parents=True, exist_ok=True)
    except Exception as e:
        raise error.ProjectInitError(e)
