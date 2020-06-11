from setuptools import setup, find_packages

with open('README.md') as f:
    readme = f.read()

with open('LICENSE') as f:
    license = f.read()

setup(
    name='Rikyu',
    version='0.1.0',
    packages=find_packages(exclude=('tests', 'docs')),
    url='https://github.com/liampulles/rikyu',
    license=license,
    author='Liam Pulles',
    author_email='me@liampulles.com',
    description='A DSL for extract and encoding DVDs',
    long_description=readme,
)
