class Error(Exception):
    pass


class ProjectInitError(Error):
    def __init__(self, cause: Exception):
        self.cause = cause


class ConfigIOError(Error):
    def __init__(self, cause: Exception):
        self.cause = cause


class HashIOError(Error):
    def __init__(self, cause: Exception):
        self.cause = cause
