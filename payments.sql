-- The DDL has three tables what is mainly used by payment processes.

CREATE TABLE users
(
	-- HubSynch ID
	-- = consumer_id
	user_id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'HubSynch ID
= consumer_id',
	email varchar(255) NOT NULL,
	-- The value is hashed by sha256
	password varchar(100) NOT NULL,
	-- NEED sender's wallet address
	last_name varchar(100) NOT NULL,
	first_name varchar(100) NOT NULL,
	-- 1978-02-20
	birthday date,
	blood tinyint unsigned,
	sex tinyint unsigned,
	-- 職業：公務員 コンサルタント コンピューター関連技術職 コンピューター関連以外の技術職 金融関係 医師 弁護士 総務・人事・事務 営業・販売 研究・開発 広報・宣伝 企画・マーケティング デザイン関係 会社経営・役員 出版・マスコミ関係 学生・フリーター 主婦 その他 不明
	-- occupation smallint unsigned DEFAULT 999 NOT NULL COMMENT '職業：公務員 コンサルタント コンピューター関連技術職 コンピューター関連以外の技術職 金融関係 医師 弁護士 総務・人事・事務 営業・販売 研究・開発 広報・宣伝 企画・マーケティング デザイン関係 会社経営・役員 出版・マスコミ関係 学生・フリーター 主婦 その他 不明',
	-- country smallint unsigned DEFAULT 1 NOT NULL,
	-- 2160033
	zip_code varchar(20),
	-- Prefecture number:
	-- http://nlftp.mlit.go.jp/ksj/gml/codelist/PrefCd.html
	address_1 smallint unsigned,
	-- 城市名稱，填 Kanji (漢字)即可，不用英文:港区
	address_2 varchar(255),
	-- 城市以外的地址：西麻布 3丁目 21-20 霞町コーポ B1F 霞町コーポ B1F → building
	address_3 varchar(255),
	-- 手機或家用電話
	tel varchar(20),
	-- 帳號停權與否
	activate_flag enum('true','false') DEFAULT 'false' NOT NULL,
	-- a mark for delete, not physical delete, but logical delete
	delete_flag enum('true','false') DEFAULT 'false' NOT NULL,
	-- 日本時區: 2014-05-09 18:41:20 
	create_timestamp datetime NOT NULL,
	-- 日本時區: 2014-05-16 22:50:42 
	update_timestamp datetime NOT NULL,
	PRIMARY KEY (user_id),
	UNIQUE (email)
) ENGINE = InnoDB;

CREATE TABLE transactions
(
	transaction_id int(10) unsigned NOT NULL AUTO_INCREMENT,
	payment_company_code smallint NOT NULL,
	-- payment_id 是 payment agent 給的  (IDHub 給的) 必須 unique 是 primary key
	-- 例如：03a9fecedcd391ac80955efb40bcbb74
	-- later 給範例：
	-- const HIVELOCITY = 999;
	-- const SMBC       = 200;
	-- const SBPS       = 100;
	-- const PAYGENT    = 10;
	-- const SMARTPIT   = 20;
	-- const GMO        = 30;
	-- const CASH       = 40;
	-- const FREE       = 50;
	-- const COLLECT    = 60;
	-- const PAYPAL    = 120;
	payment_id varchar(255) NOT NULL,
	-- later 給範例：
	-- const MOUSHIKOMI                  = 10;  // 申込済
	-- const SOKUHOU                     = 39;  // 速報
	-- const KESHIKOMI                   = 40;  // 消込済
	-- const KESHIKOMI_UNDER             = 41;  // 消込済（一部入金）
	-- const KESHIKOMI_OVER              = 42;  // 消込済（過入金）
	-- const KESHIKOMI_REQUEST           = 43;  // 消込済（依頼入金）
	-- const KESHIKOMI_REQUEST_OVER_LAP  = 44;  // 消込済（依頼入金・重複）
	-- const KESHIKOMI_COMPLETE          = 45;  // 消込済（対応済）
	-- const CANCEL                      = 60;  // 売上取消済
	-- const SOKUHO_CANCEL               = 61;  // 速報取消済
	-- const PAID_CANCEL                 = 64;  // 支払済キャンセル
	-- const CANCEL_COMPLETE             = 65;  // 支払済キャンセル（対応済）
	-- const FAILURE                     = 80;  // 決済失敗
	-- const PAY_OVER                    = 90;  // 支払期限超過
	payment_status varchar(15) NOT NULL,
	-- 價格 + 稅，不含運費 整數，JPY 沒有小數
	base_price_amount int unsigned,
	-- 給 payment agent 的錢，例如便利商店
	commission int unsigned DEFAULT 0 NOT NULL,
	-- 運費
	transportation_fee int DEFAULT 0 NOT NULL,
	-- 全部加起來，buyer 實際付的錢
	-- 例：原價：10,000、tax 800、base_price_amount、￥10,800
	--  commission ￥210、 transportation_fee ￥800。payment_amount = ￥11,810
	payment_amount int unsigned DEFAULT 0 NOT NULL,
	last_name varchar(100) DEFAULT NULL,
	first_name varchar(100) DEFAULT NULL,
	tel varchar(20),
	-- payment deadline   2014-07-29 23:59:59
	payment_limit_date date,
	create_timestamp datetime,
	update_timestamp datetime NOT NULL,
	PRIMARY KEY (use_app_transaction_id)
) ENGINE = InnoDB;
