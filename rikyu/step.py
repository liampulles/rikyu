import shutil
from pathlib import Path
from typing import List, Callable

from rikyu import error, file, hash


class StepRecord:
    def __init__(self, op_id: str, dir_hash: str):
        self.op_id = op_id
        self.dir_hash = dir_hash


class StepDefinition:
    def __init__(self, op_id: str, path: Path, step_exec: Callable[[], None]):
        self.op_id = op_id
        self.path = path
        self.step_exec = step_exec
        self.should_skip = False


def generate_op_id(op_name: str, op_in: str, op_out: str) -> str:
    return f"{op_name} $ {op_in} $ {op_out}"


def run_steps(steps: List[StepDefinition]):
    # Init steps
    for step in steps:
        already_ran = init(step.path, step.op_id)
        step.should_skip = already_ran
    # Run the steps
    for step in steps:
        if step.should_skip:
            print(f"-> Skipping step (already done): {step.op_id}")
            continue
        # Create init record
        rikyu_file = Path(step.path, ".rikyu.step")
        file.save_obj(StepRecord(step.op_id, hash.dir_hash(step.path, hash.exclude_rikyu_pattern)), rikyu_file)
        print(f"-> Running step: {step.op_id}")
        step.step_exec()
        # Update record with final dir hash
        file.save_obj(StepRecord(step.op_id, hash.dir_hash(step.path, hash.exclude_rikyu_pattern)), rikyu_file)


def init(step_path: Path, op_id: str) -> bool:
    # Create dir if not exists, or fail if not dir
    if not step_path.exists():
        step_path.mkdir(parents=True)
    if not step_path.is_dir():
        raise error.StepIOError(message=str(step_path) + " is not a directory")

    # If empty or non-existent -> proceed
    if len(list(step_path.iterdir())) == 0:
        return False

    # If not-empty, but no .rikyu file exists -> fail (not managed by rikyu)
    rikyu_step_path = Path(step_path, ".rikyu.step")
    if not rikyu_step_path.exists():
        raise error.StepIOError(message=str(step_path) + "is a non-empty directory not managed by Rikyu. Please "
                                        + "delete the directory or use another name.")
    step_record = file.load_obj(rikyu_step_path)

    # Previous hashes do not match current, consider partial -> empty dir and proceed
    current_hash = hash.dir_hash(step_path, hash.exclude_rikyu_pattern)
    if step_record.dir_hash != current_hash:
        _clear_dir(step_path)
        return False

    # Previous hashes do match, but Op identifiers do not -> fail (should use different dir)
    if step_record.op_id != op_id:
        raise error.StepConflictError(message="Op identifier mismatch: A different step finished it's task in this "
                                              + "directory - cannot reuse. Please specify another directory or delete "
                                              + "this one.")

    # Otherwise, skip
    return True


def _clear_dir(path: Path):
    print(f"(Removing partially complete job at {path})")
    for f in path.iterdir():
        if f.is_file():
            f.unlink(missing_ok=True)
        else:
            shutil.rmtree(f)
