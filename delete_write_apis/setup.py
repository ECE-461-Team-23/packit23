"""setup file for flask app"""
from setuptools import setup

setup(
    name='delete_write_apis',
    packages=['src'],
    include_package_data=True,
    install_requires=[
        'flask',
    ],
)