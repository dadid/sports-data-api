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
import pandas as pd
import numpy as np
import psycopg2
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
        self.selenium_data_queue = Queue()
        self.worker_queue = Queue()
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

    def init_selenium_workers(self):
        self._selenium_workers = {}
        if self._driver_opts is not None:
            options = webdriver.ChromeOptions()
            for opt in self._driver_opts:
                options.add_argument(opt)
        for worker_id in range(self.num_threads):
            if self._driver_opts is not None:
                self._selenium_workers[worker_id] = webdriver.Chrome(options=options)
            else:    
                self._selenium_workers[worker_id] = webdriver.Chrome()
            self.worker_queue.put(worker_id)
    
    def selenium_task(self, worker: webdriver, data):
        teamname = data
        for _, value in self.data_dict.items(): # loop over all URL sets in dictionary
            try:
                worker.get(value["url"].format(teamname))
            except TimeoutException:
                time.sleep(3)
                try:
                    worker.get(value["url"].format(teamname))
                except TimeoutException:
                    insert_audit(teamname, 2, error='Timeout on get request.')
                    continue
            for index, tag in value["html_tags"].items(): # loop over HTML tags associated with URL
                try:
                    webelem = WebDriverWait(worker, 40).until(
                        ec.presence_of_element_located(
                            (By.XPATH, tag)))
                except TimeoutException:
                    insert_audit(teamname, 3, error=f'html table - {tag}')
                    continue
                try:
                    df = html_to_dataframe(webelem, teamname)
                except ValueError:
                    insert_audit(teamname, 4, error=f'html table - {tag}')
                insert_data(df, index, teamname)
            time.sleep(2 ** np.random.randint(3, 5))

    def selenium_queue_listener(self, data_queue: Queue, worker_queue: Queue):
        logger.info('selenium listener started')
        while True:
            current_data = data_queue.get()
            if current_data == '_kill_':
                logger.warning('_kill_ encountered')
                data_queue.put(current_data)
                break
            logger.info(f'Pulled {current_data[0]} from queue')
            worker_id = worker_queue.get()
            worker = self._selenium_workers[worker_id] 
            self.selenium_task(worker, current_data)
            worker_queue.put(worker_id)

    def run(self):
        self.init_selenium_workers()
        logger.info('starting selenium processes')
        selenium_processes = [Thread(target=self.selenium_queue_listener, args=(self.selenium_data_queue, self.worker_queue)) for i in self.num_threads]
        for p in selenium_processes:
            p.daemon = True
            p.start()
        data = create_data_list()
        for d in data:
            self.selenium_data_queue.put(d)
        self.selenium_data_queue.put('_kill_')
        for p in selenium_processes:
            p.join()
        for worker in self._selenium_workers.values():
            worker.quit()

def init_db_conn():
    try:
        conn = psycopg2.connect(
            user=os.environ['SBD_DB_USER'],
            password=os.environ['SBD_DB_PASSWORD'],
            host=os.environ['SBD_DB_HOST'],
            port=os.environ['SBD_DB_POST'],
            database=os.environ['SBD_DATABASE'])
        return conn
    except psycopg2.Error as e:
        raise ValueError(f'Connection to database failed! - {e.pgerror}.')

def html_to_dataframe(webelem, teamname: str) -> pd.DataFrame:
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
        df = df[df["Rk"] != "Rk"]
    if "Name" in df.columns:
        df = df[df["Name"] != "Name"]
        df = df[~df["Name"].str.contains("total", case=False)]
        df = df[~df["Name"].str.contains("rank in 15", case=False)]
        df = df[~df["Name"].str.contains("average", case=False)]
    df["Team"] = teamname
    df = df[['Team'] + [col for col in df.columns if col != 'Team']]

    return df

def create_data_list() -> list:
    conn = init_db_conn()
    data_list = []
    with conn.cursor() as cur:
        cur.execute('SELECT teamabbrev FROM baseballreference.team WHERE ORDER BY id')
        for row in cur:
            data_list.append(row[0])
    conn.close()
    
    return data_list

def insert_audit(teamname: str, statusid: int, index: int=None, error: str=None, conn=None):
    if conn is None:
        conn = init_db_conn()
    with conn.cursor() as cur:
        try:
            cur.execute(AUDIT_INSERT, {
                "statusid": statusid, 
                "teamname": teamname,
                "tablename": index if index is None else MASTER_DICT[index]["tablename"], 
                "error": error})
            conn.commit()
        except psycopg2.Error as e:
            logger.info(e)
            exit(1)
        finally:
            conn.close()

def insert_data(df, index: int, teamname: str, conn=None):
    if conn is None:
        conn = init_db_conn()
    df_tuples = tuple(tuple(x) for x in df.values) # convert dataframe to a tuple of tuples; each inner tuple is one row
    with conn.cursor() as cur:
        values_string = ','.join(cur.mogrify(MASTER_DICT[index]["insertvalues"], x).decode('utf-8') for x in df_tuples)
        values_string = values_string.replace("'NaN'::float", 'NULL')
        query = MASTER_DICT[index]["insertquery"].format(values_string)
        try:
            cur.execute(query)
            conn.commit()
            insert_audit(teamname, 0, index=index, conn=conn)
        except psycopg2.Error as e:
            insert_audit(teamname, 1, index=index, error=e, conn=conn)
        finally: 
            conn.close()

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
    bot = SeleniumCrawler(
        driver_opts=driver_opts, 
        num_threads=5,
        selenium_data=create_data_list())
    bot.run()
    docker_exec_database_backup(zipfile=True)

if __name__ == '__main__':
    main()