import logging, time, datetime
from pathlib import Path
from master_dict import *
from multiprocessing import Queue
from threading import Thread
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as ec
from selenium.common.exceptions import NoSuchElementException, TimeoutException
from concurrent.futures import ThreadPoolExecutor
import pandas as pd
import numpy as np
import psycopg2
from sqlalchemy import create_engine, sql, bindparam
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.engine import url
import subprocess
import os

logger = logging.getLogger(__name__)
handler = logging.StreamHandler()
formatter = logging.Formatter('%(asctime)s - %(name)s:%(funcName)s:%(lineno)d - %(levelname)s - %(message)s')
handler.setFormatter(formatter)
logger.addHandler(handler)
logger.setLevel(logging.INFO)

class SeleniumCrawler:

    def __init__(self, driver_opts: list=None, num_threads: int=1):
        self._driver_opts = driver_opts
        self.num_threads = num_threads
        self.data_queue = Queue()
        self.worker_queue = Queue()
        self.database = Database(
            driver='postgresql+psycopg2',
            user=os.environ['SBD_DB_USER'],
            password=os.environ['SBD_DB_PASSWORD'],
            host=os.environ['SBD_DB_HOST'],
            port=os.environ['SBD_DB_PORT'],
            database=os.environ['SBD_DATABASE']
        )
        self.data_dict = {
            "batting": {
                "url": MASTER_DICT[0]["url"],
                "html_tags": {
                    0: MASTER_DICT[0]["htmltag"],
                    1: MASTER_DICT[1]["htmltag"]
                }
            },
            "pitching": {
                "url": MASTER_DICT[2]["url"],
                "html_tags": {
                    2: MASTER_DICT[2]["htmltag"],
                    3: MASTER_DICT[3]["htmltag"]
                }
            },
            "batting_splits": {
                "url": MASTER_DICT[4]["url"],
                "html_tags": {
                    4: MASTER_DICT[4]["htmltag"],
                    5: MASTER_DICT[5]["htmltag"]
                }
            },
            "pitching_splits": {
                "url": MASTER_DICT[6]["url"],
                "html_tags": {
                    6: MASTER_DICT[6]["htmltag"],
                    7: MASTER_DICT[7]["htmltag"]
                }
            }
        }

    def init_workers(self):
        self.workers = {}
        options = None
        if self._driver_opts is not None:
            options = webdriver.ChromeOptions()
            for opt in self._driver_opts:
                options.add_argument(opt)
        for worker_id in range(self.num_threads):
            self.workers[worker_id] = webdriver.Chrome(options=options)
            self.worker_queue.put(worker_id)
    
    def task(self, data):
        teamid = data[0]
        teamname = data[1]
        worker_id = self.worker_queue.get()
        worker = self.workers[worker_id] 
        for key, value in self.data_dict.items():
            try:
                worker.get(value["url"].format(teamname))
            except TimeoutException:
                time.sleep(3)
                try:
                    worker.get(value["url"].format(teamname))
                except TimeoutException:
                    self.insert_audit(teamid, 2, error='Timeout on get request.')
            for index, tag in value["html_tags"].items(): # loop over HTML tags associated with URL
                try:
                    webelem = WebDriverWait(worker, 40).until(
                        ec.presence_of_element_located(
                            (By.XPATH, tag)))
                except TimeoutException:
                    self.insert_audit(teamid, 3, error=f'html table - {tag}')
                    continue
                try:
                    df = self.html_to_dataframe(webelem, teamid)
                except ValueError:
                    self.insert_audit(teamid, 4, error=f'html table - {tag}')
                self.insert_data(df, index, teamid)
            time.sleep(2 ** np.random.randint(3, 5))
        self.worker_queue.put(worker_id)
    
    def html_to_dataframe(self, webelem, teamid) -> pd.DataFrame:
        try:
            html = webelem.get_attribute('outerHTML')
            df = pd.read_html(html)[0]
        except TimeoutError:
            try:
                html = webelem.get_attribute('outerHTML')
                df = pd.read_html(html)[0]
            except:
                raise ValueError('Error getting outerHTML attribute')
        if "Rk" in df.columns:
            df = df[df.Rk != "Rk"]
        if "Name" in df.columns:
            df = df[df.Name != "Name"]
            df = df[~df.Name.str.contains("total", case=False)]
            df = df[~df.Name.str.contains("rank in 15", case=False)]
            df = df[~df.Name.str.contains("average", case=False)]
        df.columns = df.columns.str.lower()
        df.rename(
            mapper={
                'w-l%': 'wl',
                'era+': 'eraplus',
                'so/w': 'sow',
                'ops+': 'opsplus',
                'tops+': 'topsplus',
                'sops+': 'sopsplus',
                'xbt%': 'xbtpct',
                'rs%': 'rspct',
                'sb%': 'sbpct',
                '2b':  'twob',
                '3b':  'threeb'
                },
            axis='columns',
            inplace=True)
        df["teamid"] = teamid

        return df
    
    def create_data_list(self):
        if not hasattr(self.database, 'engine'):
            self.database.connect()
        with self.database.engine.connect() as conn:
            res = conn.execute(sql.text('SELECT id, teamabbrev FROM baseballreference.team ORDER BY id'))
            self.data = [(row['id'], row['teamabbrev']) for row in res]
    
    def insert_audit(self, teamid, statusid, index=None, error=None):
        if not hasattr(self.database, 'engine'):
            self.database.connect()
        with self.database.engine.connect() as conn:
            conn.execute(sql.text(AUDIT_INSERT), {
                "statusid": statusid, 
                "teamid": teamid,
                "tablename": index if index is None else MASTER_DICT[index]["tablename"], 
                "error": str(error) if error is not None else error})

    def insert_data(self, df, index, teamid):
        if not hasattr(self.database, 'engine'):
            self.database.connect()
        try:
            df.to_sql(
                name=MASTER_DICT[index]["table"],
                schema='baseballreference',
                con=self.database.engine,
                index=False,
                if_exists='append',
                chunksize=1000
            )
            self.insert_audit(teamid, 0, index=index)
        except SQLAlchemyError as e:
            self.insert_audit(teamid, 1, index=index, error=e)
    
    def listen_depr(self, data_queue: Queue, worker_queue: Queue):
        logger.info('listener started')
        while True:
            current_data = data_queue.get()
            if current_data == '_kill_':
                logger.warning('_kill_ found')
                data_queue.put(current_data)
                break
            logger.info(f'Pulled {current_data} from queue')
            worker_id = worker_queue.get()
            worker = self.workers[worker_id] 
            self.task(current_data, worker)
            worker_queue.put(worker_id)

    def run_depr(self):
        self.init_workers()
        self.create_data_list()
        logger.info('starting workers')
        selenium_processes = [Thread(target=self.listen, args=(self.data_queue, self.worker_queue)) for i in range(self.num_threads)]
        for p in selenium_processes:
            p.daemon = True
            p.start()
        for d in self.data:
            self.data_queue.put(d)
        self.data_queue.put('_kill_')
        for p in selenium_processes:
            p.join()
        for worker in self.workers.values():
            worker.quit()

    def run(self):
        self.create_data_list()
        self.init_workers()
        with ThreadPoolExecutor(max_workers=self.num_threads) as executor:
            executor.map(self.task, self.data)


class Database:
    
    def __init__(self, driver, host, database, user, password, port):
        self.driver = driver
        self.host = host
        self.database = database
        self.user = user
        self.password = password
        self.port = port

    def connect(self, execute_many=False):
        conn_str = url.URL(
            drivername=self.driver,
            username=self.user,
            password=self.password,
            host=self.host,
            port=self.port,
            database=self.database
        )
        if self.driver == 'mssql+pyodbc':
            self.engine = create_engine(conn_str, fast_executemany=execute_many)
        elif self.driver == 'postgresql+psycopg2': 
            self.engine = create_engine(conn_str)


def docker_exec_database_backup(zipfile=False):
    backup_path = Path(r'C:\Users\Daniel\01.devel\sports-data-api\db-backups')
    getdate = datetime.datetime.now().strftime('%Y-%m-%d_%H%M%S')
    if zipfile is False:
        backup_file = backup_path / f'baseball_ref_db_backup_{getdate}.sql'
        subprocess.call(f'docker exec -t dev-postgres pg_dumpall --no-owner -c -U postgres | {backup_file}', shell=True)
    else:
        backup_file = backup_path / f'baseball_ref_db_backup_{getdate}.zip'
        subprocess.call(f'docker exec -t dev-postgres pg_dumpall --no-owner -c -U postgres | gzip > {backup_file}', shell=True)


def main():
    user_agent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36'
    ext_path = Path.home() / '01.devel' / 'chrome_anti_detection_extension'
    driver_opts = [f'user-agent={user_agent}', 'log-level=3'] # f'load-extension={ext_path.absolute()}']
    crawler = SeleniumCrawler(driver_opts=driver_opts, num_threads=4)
    crawler.run()


if __name__ == '__main__':
    main()