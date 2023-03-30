package repository

// DDL
const (
	ddl = `
	create table sois (
		id integer not null primary key autoincrement,
		bucket text not null,
		hash text not null,
		name text,
		title text,
		path text not null, 
		uri text,
		rate real,
		num_views integer,
		num_reads integer,
		comprehension integer,
		created_at datetime
	);
	create index idx_soi_hash on sois(hash);
	create index idx_soi_path on sois(path);
	create index idx_soi_name on sois(name);

	create table ogs (
		id integer not null primary key autoincrement,
		soi_id integer not null,
		title text,
		url text,
		type text,
		description text,
		site_name text
	);

	create table og_imgs (
		id integer not null primary key autoincrement,
		og_id integer not null,
		url text,
		secure_url text,
		type text,
		width integer,
		height integer,
		alt text
	);

	create table usage_logs (
		soi_id integer not null,
		type text,
		used_at datetime
	);

	create table tags (
		id integer not null primary key autoincrement,
		name text
	);

	create table soi_tags (
		soi_id integer not null,
		tag_id integer not null
	);

	create table kv_tags (
		id integer not null primary key autoincrement,
		key text not null,
		value text not null
	);

	create table soi_kv_tags (
		soi_id integer not null,
		kv_tag_id integer not null
	);`
)

var (
	// sois
	selectAllSoisQuery = `
		select 
			id, hash, name, title, path, uri, rate, num_views, num_reads, comprehension, created_at 
		from
			sois
	where
	   bucket = ? `
	selectSoiQuery       = selectAllSoisQuery + " and hash = ? order by path, name desc "
	selectSoiByPathQuery = selectAllSoisQuery + " and path = ? order by path, name desc "
	storeSoiQuery        = `
		insert into sois 
			(bucket, hash, name, title, path, uri, rate, num_views, num_reads, comprehension, created_at) 
		values 
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) `

	//selectPathQuery = "select path, name, path || '/' || name as full_path from sois where bucket = ? and path like ? order by full_path desc"

	// tags
	selectTagQuery    = "select name from tags t join soi_tags st on (t.id == st.tag_id) where st.soi_id = ?"
	storeTagQuery     = "insert into tags (name) values (?)"
	storeSoiTagsQuery = "insert into soi_tags (soi_id, tag_id) values (?, ?)"

	// kv_tags
	selectKVTagQuery    = "select key, value from kv_tags kvt join soi_kv_tags skvt on (kvt.id == skvt.kv_tag_id) where skvt.soi_id = ?"
	storeKVTagQuery     = "insert into kv_tags (key, value) values (?, ?)"
	storeSoiKVTagsQuery = "insert into soi_kv_tags (soi_id, kv_tag_id) values (?, ?)"

	// usage_logs
	selectUsageLogQuery = "select soi_id, type, used_at from usage_logs where soi_id = ? order by used_at desc"
	storeUsageLogQuery  = "insert into usage_logs (soi_id, type, used_at) values (?, ?, ?)"

	// ogs
	selectOGQuery = "select id, soi_id, title, url, type, description, site_name from ogs where soi_id = ?"
	storeOGQuery  = "insert into ogs (soi_id, title, url, type, description, site_name) values (?, ?, ?, ?, ?, ?)"

	// og_imgs
	selectOGImgQuery = "select id, og_id, url, secure_url, type, width, height, alt from og_imgs where og_id = ?"
	storeOGImgQuery  = "insert into og_imgs (og_id, url, secure_url, type, width, height, alt) values (?, ?, ?, ?, ?, ?, ?)"
)
