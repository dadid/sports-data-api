MASTER_DICT = {
    0: {
        "tablename": 'batting',
        "url": 'https://www.baseball-reference.com/teams/{}/{}-batting.shtml#team_batting::none',
        "htmltag": '//*[@id="team_batting"]',
    },
    1: {
        "tablename": 'baserunning',
        "url": 'https://www.baseball-reference.com/teams/{}/{}-batting.shtml#team_batting::none', 
        "htmltag": '//*[@id="players_baserunning_batting"]',
    },
    2: {
        "tablename": 'pitching',
        "url": 'https://www.baseball-reference.com/teams/{}/{}-pitching.shtml',
        "htmltag": '//*[@id="team_pitching"]',
    },
    3: {
        "tablename": 'batting_pitching',
        "url": 'https://www.baseball-reference.com/teams/{}/{}-pitching.shtml',
        "htmltag": '//*[@id="players_batting_pitching"]',
    },
    4: {
        "tablename": 'batting_splits',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=b&team={}&year={}',
        "htmltag": '//*[@id="plato"]',
    },
    5: {
        "tablename": 'batting_home_away',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=b&team={}&year={}',
        "htmltag": '//*[@id="hmvis"]',
    },
    6: {
        "tablename": 'pitching_splits',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=p&team={}&year={}',
        "htmltag": '//*[@id="plato"]',
    },
    7: {
        "tablename": 'pitching_home_away',
        "url": 'https://www.baseball-reference.com/teams/split.cgi?t=p&team={}&year={}',
        "htmltag": '//*[@id="hmvis"]',
    }
}

AUDIT_INSERT = "INSERT INTO baseballreference.audit VALUES (default, :statusid, :teamid, :tablename, :error, default)"