create table prices(
	id integer primary key autoincrement,
  	code text,
  	codein text,
  	name text,
  	high decimal(10,5),
  	low decimal(10,5),
  	varBid decimal(10,5),
  	pctChange decimal(10,5),
  	bid decimal(10,5),
  	ask decimal(10,5),
  	timestamp integer,
  	create_date text
);