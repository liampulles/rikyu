from .context import error, hash

import os
import unittest
from pathlib import Path
from typing import List

THIS_DIR = os.path.dirname(os.path.abspath(__file__))


class FileHashTestSuite(unittest.TestCase):
    def test_file_hash__given_invalid_file__should_fail(self):
        # Setup fixture
        fixture_base = _testdata_path(["hash"])
        fixture_path = _testdata_path(["hash", "does_not_exist"])

        # Exercise SUT
        def sut():
            hash.file_hash(fixture_base, fixture_path)

        # Verify results
        self.assertRaises(error.HashIOError, sut)

    def test_file_hash__given_two_empty_files_with_same_relative_path__should_have_same_hash(self):
        # Setup fixture
        fixture_base1 = _testdata_path(["hash", "empty1"])
        fixture_path1 = _testdata_path(["hash", "empty1", "somedir", "empty.txt"])
        fixture_base2 = _testdata_path(["hash", "empty2"])
        fixture_path2 = _testdata_path(["hash", "empty2", "somedir", "empty.txt"])

        # Exercise SUT
        actual1 = hash.file_hash(fixture_base1, fixture_path1)
        actual2 = hash.file_hash(fixture_base2, fixture_path2)

        # Verify results
        self.assertEqual(actual1, actual2)

    def test_file_hash__given_two_empty_files_with_different_names__should_have_different_hash(self):
        # Setup fixture
        fixture_base1 = _testdata_path(["hash", "empty1"])
        fixture_path1 = _testdata_path(["hash", "empty1", "somedir", "empty.txt"])
        fixture_base2 = _testdata_path(["hash", "empty1"])
        fixture_path2 = _testdata_path(["hash", "empty1", "somedir", "another_empty.txt"])

        # Exercise SUT
        actual1 = hash.file_hash(fixture_base1, fixture_path1)
        actual2 = hash.file_hash(fixture_base2, fixture_path2)

        # Verify results
        self.assertNotEqual(actual1, actual2)

    def test_file_hash__given_two_empty_files_with_different_relative_paths__should_have_different_hash(self):
        # Setup fixture
        fixture_base1 = _testdata_path(["hash", "empty1"])
        fixture_path1 = _testdata_path(["hash", "empty1", "somedir", "empty.txt"])
        fixture_base2 = _testdata_path(["hash", "empty1"])
        fixture_path2 = _testdata_path(["hash", "empty1", "anotherdir", "empty.txt"])

        # Exercise SUT
        actual1 = hash.file_hash(fixture_base1, fixture_path1)
        actual2 = hash.file_hash(fixture_base2, fixture_path2)

        # Verify results
        self.assertNotEqual(actual1, actual2)

    def test_file_hash__given_two_same_files_with_same_path__should_have_same_hash(self):
        # Setup fixture
        fixture_base1 = _testdata_path(["hash", "normal1"])
        fixture_path1 = _testdata_path(["hash", "normal1", "somedir", "file.txt"])
        fixture_base2 = _testdata_path(["hash", "normal2"])
        fixture_path2 = _testdata_path(["hash", "normal2", "somedir", "file.txt"])

        # Exercise SUT
        actual1 = hash.file_hash(fixture_base1, fixture_path1)
        actual2 = hash.file_hash(fixture_base2, fixture_path2)

        # Verify results
        self.assertEqual(actual1, actual2)

    def test_file_hash__given_two_different_files_with_same_path__should_have_different_hash(self):
        # Setup fixture
        fixture_base1 = _testdata_path(["hash", "normal1"])
        fixture_path1 = _testdata_path(["hash", "normal1", "somedir", "file.txt"])
        fixture_base2 = _testdata_path(["hash", "normal3"])
        fixture_path2 = _testdata_path(["hash", "normal3", "somedir", "file.txt"])

        # Exercise SUT
        actual1 = hash.file_hash(fixture_base1, fixture_path1)
        actual2 = hash.file_hash(fixture_base2, fixture_path2)

        # Verify results
        self.assertNotEqual(actual1, actual2)


class DirHashTestSuite(unittest.TestCase):
    def test_dir_hash__given_two_empty_trees__should_have_same_hash(self):
        # Setup fixture
        fixture1 = _testdata_path(["hash", "emptydir1"])
        fixture2 = _testdata_path(["hash", "emptydir2"])

        # Exercise SUT
        actual1 = hash.dir_hash(fixture1)
        actual2 = hash.dir_hash(fixture2)

        # Verify results
        self.assertEqual(actual1, actual2)

    def test_dir_hash__given_two_different_trees_with_empty_folders__should_have_same_hash(self):
        # Setup fixture
        fixture1 = _testdata_path(["hash", "emptyfolders1"])
        fixture2 = _testdata_path(["hash", "emptyfolders2"])

        # Exercise SUT
        actual1 = hash.dir_hash(fixture1)
        actual2 = hash.dir_hash(fixture2)

        # Verify results
        self.assertEqual(actual1, actual2)

    def test_dir_hash__given_two_same_trees_with_same_files__should_have_same_hash(self):
        # Setup fixture
        fixture1 = _testdata_path(["hash", "samefilesdir1"])
        fixture2 = _testdata_path(["hash", "samefilesdir2"])

        # Exercise SUT
        actual1 = hash.dir_hash(fixture1)
        actual2 = hash.dir_hash(fixture2)

        # Verify results
        self.assertEqual(actual1, actual2)

    def test_dir_hash__given_two_same_trees_with_different_files__should_have_different_hash(self):
        # Setup fixture
        fixture1 = _testdata_path(["hash", "samefilesdir1"])
        fixture2 = _testdata_path(["hash", "samefilesdir3"])

        # Exercise SUT
        actual1 = hash.dir_hash(fixture1)
        actual2 = hash.dir_hash(fixture2)

        # Verify results
        self.assertNotEqual(actual1, actual2)


def _testdata_path(locations: List[str]) -> Path:
    base = Path(THIS_DIR, "testdata")
    for location in locations:
        base = Path(base, location)
    return base


if __name__ == '__main__':
    unittest.main()
