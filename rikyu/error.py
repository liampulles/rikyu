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


class RipperIOError(Error):
    def __init__(self, cause: Exception = None, message: str = None):
        if cause:
            self.cause = cause
        if message:
            self.message = message


class StepIOError(Error):
    def __init__(self, cause: Exception = None, message: str = None):
        if cause:
            self.cause = cause
        if message:
            self.message = message


class StepConflictError(Error):
    def __init__(self, cause: Exception = None, message: str = None):
        if cause:
            self.cause = cause
        if message:
            self.message = message
