create database tokenizer;
use tokenizer;
create table tokenizer_metadata (request_id varchar(255), global_tx_id varchar(255), created_date timestamp, language varchar(255), PRIMARY KEY (request_id));