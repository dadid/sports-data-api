import logging, time, datetime, glob
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
from sqlalchemy import create_engine, sql
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.engine import url
from sqlalchemy import bindparam
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
        # self.selenium_data_queue = Queue()
        # self.worker_queue = Queue()
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
    
    def selenium_task(self, data):
        teamname = data
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
                    # self.insert_audit(teamname, 2, error='Timeout on get request.')
                    continue
            # for index, tag in value["html_tags"].items(): # loop over HTML tags associated with URL
            #     try:
            #         webelem = WebDriverWait(worker, 40).until(
            #             ec.presence_of_element_located(
            #                 (By.XPATH, tag)))
            #     except TimeoutException:
            #         self.insert_audit(teamname, 3, error=f'html table - {tag}')
            #         continue
            #     try:
            #         df = self.html_to_dataframe(webelem, teamname)
            #     except ValueError:
            #         self.insert_audit(teamname, 4, error=f'html table - {tag}')
            #     self.insert_data(df, index, teamname)
            # time.sleep(2 ** np.random.randint(3, 5))
        self.worker_queue.put(worker_id)
    
    def html_to_dataframe(self, webelem, teamname: str) -> pd.DataFrame:
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
        df["Team"] = teamname
        df = df[['Team'] + [col for col in df.columns if col != 'Team']]

        return df
    
    def create_data_list(self):
        with self.database.connect() as conn:
            self.data = []
            res = conn.execute(sql.text('SELECT teamabbrev FROM baseballreference.team WHERE ORDER BY id'))
            for row in res:
                self.data.append(row['teamabbrev'])
    
    def insert_audit(self, teamname, statusid, index=None, error=None):
        with self.database.connect() as conn:
            conn.execute(AUDIT_INSERT, {
                "statusid": statusid, 
                "teamname": teamname,
                "tablename": index if index is None else MASTER_DICT[index]["tablename"], 
                "error": error})

    def insert_data(self, df, index, teamname):
        with self.database.connect() as conn:
            try:
                df.to_sql(
                    name=MASTER_DICT[index]["table"],
                    con=engine,
                    index=False,
                    if_exists='append',
                    chunksize=1000
                )
                self.insert_audit(teamname, 0, index=index)
            except SQLAlchemyError as e:
                self.insert_audit(teamname, 1, index=index, error=e)

    def insert_data_depr(self, df, index, teamname):
        df_tuples = tuple(tuple(x) for x in df.values) # convert dataframe to a tuple of tuples; each inner tuple is one row
        values_string = ','.join(cur.mogrify(MASTER_DICT[index]["insertvalues"], x).decode('utf-8') for x in df_tuples)
        values_string = values_string.replace("'NaN'::float", 'NULL')
        query = MASTER_DICT[index]["insertquery"].format(values_string)
        with self.database.connect() as conn:
            try:
                conn.execute(query)
                self.insert_audit(teamname, 0, index=index)
            except SQLAlchemyError as e:
                self.insert_audit(teamname, 1, index=index, error=e)
    
    def listener(self, data_queue: Queue, worker_queue: Queue):
        logger.info('listener started')
        while True:
            current_data = data_queue.get()
            if current_data == '_kill_':
                logger.warning('_kill_ found')
                data_queue.put(current_data)
                break
            logger.info(f'Pulled {current_data[0]} from queue')
            worker_id = worker_queue.get()
            worker = self.workers[worker_id] 
            self.selenium_task(worker, current_data)
            worker_queue.put(worker_id)

    def run(self):
        self.init_workers()
        self.create_data_list()
        logger.info('starting workers')
        selenium_processes = [Thread(target=self.listener, args=(self.selenium_data_queue, self.worker_queue)) for i in self.num_threads]
        for p in selenium_processes:
            p.daemon = True
            p.start()
        for d in self.data:
            self.selenium_data_queue.put(d)
        self.selenium_data_queue.put('_kill_')
        for p in selenium_processes:
            p.join()
        for worker in self.workers.values():
            worker.quit()

    def run_threadpool(self):
        self.create_data_list()
        self.init_workers()
        with ThreadPoolExecutor(max_workers=self.num_threads) as executor:
            executor.map(self.selenium_task, self.data)


class Database:
    
    def __init__(self, driver, host, database, user, password, port, is_trusted=False, read_only=False):
        self.driver = driver
        self.host = host
        self.database = database
        self.user = user
        self.password = password
        self.port = port
        self.is_trusted = 'yes' if is_trusted else 'no'
        self.read_only = read_only

    def connect(self, execute_many=False):
        conn_str = url.URL(
            username=self.user,
            password=self.password,
            host=self.host,
            port=self.port,
            database=self.database,
            query={'driver': self.driver, 'readonly': self.read_only, 'trusted_connection': self.is_trusted}
        )
        return create_engine(conn_str, fast_executemany=execute_many)


def docker_exec_database_backup(zipfile=False):
    backup_path = Path(r'C:\Users\Daniel\01.devel\sports-data-api\db-backups')
    getdate = datetime.datetime.now().strftime('%Y-%m-%d_%H%M%S')
    if zipfile is False:
        backup_file = backup_path / f'baseball_ref_db_backup_{getdate}.sql'
        subprocess.call(f'docker exec -t dev-postgres pg_dumpall -c -U postgres | {backup_file}', shell=True)
    else:
        backup_file = backup_path / f'baseball_ref_db_backup_{getdate}.zip'
        subprocess.call(f'docker exec -t dev-postgres pg_dumpall -c -U postgres | gzip > {backup_file}', shell=True)

def main():
    user_agent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36'
    extension_dir = r'C:\Users\Daniel\01.devel\chrome_anti_detection_extension'
    driver_opts = [f'user-agent={user_agent}', 'log-level=3', f'load-extension={extension_dir}']
    bot = SeleniumCrawler(driver_opts=driver_opts, num_threads=5)
    # bot.run()
    bot.run_threadpool()
    # docker_exec_database_backup(zipfile=True)

if __name__ == '__main__':
    main()