ALTER TABLE public.users ADD password_last_updated timestamptz(6) NULL;
INSERT INTO public.sys_function (id, "name", enable_flag, description, remark) VALUES('function_4', 'shared_account', 'enable', '共享账号', NULL);

ALTER TABLE "public"."action_info"
ALTER COLUMN "action_param_type" TYPE varchar(20) COLLATE "pg_catalog"."default";