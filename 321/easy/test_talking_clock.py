import unittest
from talking_clock import translate

class TalkingClockTranslateTest(unittest.TestCase):
    def test_whole(self):
        self.assertEqual(translate('23:59'), "It's eleven fifty nine pm")
        self.assertEqual(translate('11:59'), "It's eleven fifty nine am")

    def test_tenth(self):
        self.assertEqual(translate('10:11'), "It's ten eleven am")
        self.assertEqual(translate('15:19'), "It's three nineteen pm")

    def test_zeros(self):
        self.assertEqual(translate('01:00'), "It's one am")
        self.assertEqual(translate('17:00'), "It's five pm")
        self.assertEqual(translate('09:20'), "It's nine twenty am")
        self.assertEqual(translate('13:30'), "It's one thirty pm")

    def test_am_pm(self):
        self.assertEqual(translate('12:00'), "It's twelve pm")
        self.assertEqual(translate('00:00'), "It's twelve am")


if __name__ == '__main__':
    unittest.main()
