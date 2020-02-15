MASTER_DICT = {
    0: {
        "tablename": 'batting',
        "url": 'https://www.baseball-reference.com/teams/{}/2019-batting.shtml#team_batting::none',
        "htmltag": '//*[@id="team_batting"]',
        "insertquery": 'INSERT INTO baseballreference.batting VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    1: {
        "tablename": 'baserunning',
        "url": 'https://www.baseball-reference.com/teams/{}/2019-batting.shtml#team_batting::none', 
        "htmltag": '//*[@id="players_baserunning_batting"]',
        "insertquery": 'INSERT INTO baseballreference.baserunning VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    2: {
        "tablename": 'pitching',
        "url": 'https://www.baseball-reference.com/teams/{}/2019-pitching.shtml',
        "htmltag": '//*[@id="team_pitching"]',
        "insertquery": 'INSERT INTO baseballreference.pitching VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    3: {
        "tablename": 'batting_pitching',
        "url": 'https://www.baseball-reference.com/teams/{}/2019-pitching.shtml',
        "htmltag": '//*[@id="players_batting_pitching"]',
        "insertquery": 'INSERT INTO baseballreference.batting_pitching VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    4: {
        "tablename": 'batting_splits',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=b&team={}&year=2019',
        "htmltag": '//*[@id="plato"]',
        "insertquery": 'INSERT INTO baseballreference.batting_splits VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    5: {
        "tablename": 'batting_home_away',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=b&team={}&year=2019',
        "htmltag": '//*[@id="hmvis"]',
        "insertquery": 'INSERT INTO baseballreference.batting_home_away VALUES {}', # done
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    6: {
        "tablename": 'pitching_splits',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=p&team={}&year=2019',
        "htmltag": '//*[@id="plato"]',
        "insertquery": 'INSERT INTO baseballreference.pitching_splits VALUES {}',
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    },
    7: {
        "tablename": 'pitching_home_away',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=p&team={}&year=2019',
        "htmltag": '//*[@id="hmvis"]',
        "insertquery": 'INSERT INTO baseballreference.pitching_home_away VALUES {}',
        "insertvalues": '(default,(select id from baseballreference.team where teamabbrev=%s),%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,default)\n'
    }
}

AUDIT_INSERT = "INSERT INTO baseballreference.audit VALUES (default, %(statusid)s, (select id from baseballreference.team where teamabbrev=%(teamname)s), %(tablename)s, %(error)s, default)"