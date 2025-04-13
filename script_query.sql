select * from orders  ;
select * from order_items ;

truncate table order_items ;
truncate table orders CASCADE;

SHOW max_connections;
SHOW config_file;
SELECT count(*) FROM pg_stat_activity;