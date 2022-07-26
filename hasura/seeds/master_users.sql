INSERT INTO "master_users" (id,name,is_admin)
VALUES
  (1,'テスト太郎',true),
  (2,'テスト次郎',false),
  (3,'テスト三郎',false)
ON CONFLICT ()
DO UPDATE SET
  "id" = EXCLUDED.id,
  "name" = EXCLUDED.name,
  "is_admin" = EXCLUDED.is_admin
;
