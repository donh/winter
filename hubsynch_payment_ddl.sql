-- The DDL has three tables what is mainly used by payment processes.

CREATE TABLE users
(
	-- HubSynch ID
	-- = consumer_id
	user_id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'HubSynch ID
= consumer_id',
	email varchar(255) NOT NULL,
	-- The value is hashed by sha256
	password varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	first_name varchar(100) NOT NULL,
	last_name_kana varchar(100),
	first_name_kana varchar(100),
	last_name_rome varchar(100),
	first_name_rome varchar(100),
	birthday date,
	blood tinyint unsigned,
	sex tinyint unsigned,
	-- 職業：公務員 コンサルタント コンピューター関連技術職 コンピューター関連以外の技術職 金融関係 医師 弁護士 総務・人事・事務 営業・販売 研究・開発 広報・宣伝 企画・マーケティング デザイン関係 会社経営・役員 出版・マスコミ関係 学生・フリーター 主婦 その他 不明
	occupation smallint unsigned DEFAULT 999 NOT NULL COMMENT '職業：公務員 コンサルタント コンピューター関連技術職 コンピューター関連以外の技術職 金融関係 医師 弁護士 総務・人事・事務 営業・販売 研究・開発 広報・宣伝 企画・マーケティング デザイン関係 会社経営・役員 出版・マスコミ関係 学生・フリーター 主婦 その他 不明',
	country smallint unsigned DEFAULT 1 NOT NULL,
	zip_code varchar(20),
	address_1 smallint unsigned,
	address_2 varchar(255),
	address_3 varchar(255),
	tel varchar(20),
	profile_image varchar(255),
	activate_code varchar(100),
	activate_flag enum('true','false') DEFAULT 'false' NOT NULL,
	delete_flag enum('true','false') DEFAULT 'false' NOT NULL,
	create_timestamp datetime NOT NULL,
	update_timestamp datetime NOT NULL,
	PRIMARY KEY (user_id),
	UNIQUE (email)
) ENGINE = InnoDB;

CREATE TABLE use_apps
(
	use_app_id int unsigned NOT NULL AUTO_INCREMENT,
	-- A foreign key of company's application to relate which user is using which application.
	company_app_id int unsigned NOT NULL,
	-- A foreign key of comapny to jump the relation in SQL JOIN.
	company_id int unsigned,
	user_id int unsigned,
	delete_flag enum('true','false') DEFAULT 'false' NOT NULL,
	create_timestamp datetime NOT NULL,
	update_timestamp datetime NOT NULL,
	PRIMARY KEY (use_app_id)
) ENGINE = InnoDB;

CREATE TABLE use_app_transactions
(
	use_app_transaction_id int unsigned NOT NULL AUTO_INCREMENT,
	use_app_id int unsigned NOT NULL,
	-- Foreign key of use_app_subscriptions table
	use_app_subscription_id int unsigned,
	-- Foreign key of subscribers table
	subscriber_id int unsigned,
	payment_company_code smallint NOT NULL,
	payment_id varchar(255) NOT NULL,
	payment_status varchar(15) NOT NULL,
	base_price_amount int unsigned,
	commission int unsigned DEFAULT 0 NOT NULL,
	transportation_fee int DEFAULT 0 NOT NULL,
	payment_amount int unsigned DEFAULT 0 NOT NULL,
	purchase_genre_code smallint NOT NULL,
	order_id varchar(255),
	delivering_address_last_name varchar(100) DEFAULT 'NULL',
	delivering_address_first_name varchar(255) DEFAULT 'NULL',
	delivering_address_last_name_kana varchar(100),
	delivering_address_first_name_kana varchar(100),
	delivering_address_last_name_rome varchar(100),
	delivering_address_first_name_rome varchar(100),
	zip_code varchar(20),
	country smallint unsigned,
	address_1 smallint unsigned,
	address_2 varchar(255),
	address_3 varchar(255),
	tel varchar(20),
	used_brand_creditcard varchar(50),
	valid_term varchar(4),
	holder_name varchar(255),
	payment_limit_date date,
	convenience_payment_information varchar(255),
	banktransfer_payment_information varchar(255),
	atm_payment_information varchar(255),
	remind_complete_status smallint(2),
	create_timestamp datetime,
	update_timestamp datetime NOT NULL,
	PRIMARY KEY (use_app_transaction_id)
) ENGINE = InnoDB;

