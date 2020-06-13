from pathlib import Path

from rikyu.source import Source


class Ripper:
    def __init__(self, dest_base_path: Path, dvd_path: Path):
        self.dest_base_path = dest_base_path
        # TODO: Verify above
        self.dvd_path = dvd_path
        # TODO: Verify above
        self.titles = []

    def title(self, title_no: int, dest: Source):
        self.titles += [{"title_no": title_no, "dest": Path(self.dest_base_path, dest.relative_path())}]
        return self

    def rip(self):
        # TODO: Handle destination cases for each title
        #  - If empty or non-existent -> proceed
        #  - If not-empty, but no .rikyu files exist -> fail (not managed by rikyu)
        #  - .rikyu has not stored hashes for directory, consider partial -> empty dir and proceed
        #  - .rikyu has hashes, but contents do not match hashes, consider partial -> empty dir and proceed
        #  - .rikyu has hashes, content matches hashes -> skip
        pass