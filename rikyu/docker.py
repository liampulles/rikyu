from pathlib import Path
from typing import List, Callable, Tuple

import docker
from docker.types import Mount

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
    step_exec = _step_exec(image, inputs, output, bash, args)
    return step.StepDefinition(
        op_id,
        output,
        step_exec
    )


def _step_exec(image: str,
               inputs: List[Tuple[str, Path]],
               output: Path,
               bash: str = None,
               args: List[str] = None) -> Callable[[], None]:
    if bash is not None:
        cmd = bash
    if args is not None:
        cmd = args
    mounts = [Mount(target=f"/in/{i[0]}", source=str(i[1]), read_only=True, type="bind") for i in inputs]
    mounts += [Mount(target="/out", source=str(output), read_only=False, type="bind")]

    def _func():
        client = docker.from_env()
        logs = client.containers.run(
            image=image,
            command=cmd,
            stdout=True, stderr=True,
            remove=True,
            mounts=mounts
        )
        print("===LOGS===")
        print(bytes(logs).decode("utf-8"))
        print("==========\n")

    return _func


def _op_name(image: str, bash: str = None, args: List[str] = None) -> str:
    if bash is not None:
        return f"Docker(image: {image}, cmd: \"{bash}\")"
    if args is not None:
        return f"Docker(image: {image}, cmd: \"{' '.join(args)}\")"
    return f"Docker(image: {image}, cmd: N/A)"


def _op_in(inputs: List[Tuple[str, Path]]) -> str:
    return ", ".join([f"(\"{t[0]}\" = \"{t[1]}\")" for t in inputs])
