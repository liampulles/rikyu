from pathlib import Path
from typing import List, Callable, Tuple

from rikyu import step


def docker_step_def(image: str,
                    inputs: List[Tuple[str, Path]],
                    output: Path,
                    base: Path,
                    bash: str = None,
                    args: List[str] = None) -> step.StepDefinition:
    op_id = step.generate_op_id(
        _op_name(image, bash, args),
        _op_in(inputs),
        str(output.relative_to(base))
    )
    step_exec = _step_exec(output)
    return step.StepDefinition(
        op_id,
        output,
        step_exec
    )


def _step_exec(output: Path) -> Callable[[], None]:
    def _func():
        with Path(output, "docker_test").open("w") as f:
            f.write("docker test!")

    return _func


def _op_name(image: str, bash: str = None, args: List[str] = None) -> str:
    if bash is not None:
        return f"Docker(image: {image}, cmd: \"{bash}\")"
    if args is not None:
        return f"Docker(image: {image}, cmd: \"{' '.join(args)}\")"
    return f"Docker(image: {image}, cmd: N/A)"


def _op_in(inputs: List[Tuple[str, Path]]) -> str:
    return ", ".join([f"(\"{t[0]}\" = \"{t[1]}\")" for t in inputs])
