import unittest

class TestConvertBoolean(unittest.TestCase):

    def setUp(self):
        self.val = 'sample'

    def test_pos(self):
        with self.subTest(key='Good Input Test'):
            expectedVal = 0
            realVal = 0
            self.assertEqual(expectedVal, realVal, "Result is incorrect")

    def test_neg(self):
        with self.subTest(key='Bad Input Test'):
            expectedVal = 1
            realVal = 2
            self.assertNotEqual(expectedVal, realVal, "Result is incorrect")
            # with self.assertRaises(ValueError):
            #     pass
                  # make error happen