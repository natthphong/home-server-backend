drop table tbl_user_role;
drop table tbl_user;

CREATE TABLE tbl_user
(
    user_id               VARCHAR(50) primary key ,
    first_name_th         VARCHAR(255),
    first_name_en         VARCHAR(255),
    mid_name_th           VARCHAR(255),
    mid_name_en           VARCHAR(255),
    last_name_th          VARCHAR(255),
    last_name_en          VARCHAR(255),
    phone                 VARCHAR(20),
    user_id_type          VARCHAR(50),
    email                 VARCHAR(100),
    nationality           VARCHAR(50),
    occupation            VARCHAR(100),
    request_ref           VARCHAR(100),
    birth_date            DATE,
    gender                CHAR(1),
    tax_id                VARCHAR(20),
    second_email          VARCHAR(100),
    occupation_other_desc VARCHAR(255),
    is_active             varchar(1) DEFAULT 'Y',
    password              VARCHAR(255),
    status                VARCHAR(50),
    account_name          VARCHAR(100),
    external_id           VARCHAR(100),
    user_details          JSON,
    in_active              VARCHAR(1) DEFAULT 'N',

    -- Pattern Fields
    create_at             TIMESTAMPTZ DEFAULT now(),
    create_by             VARCHAR(50),
    update_by             VARCHAR(50),
    update_at             TIMESTAMPTZ,
    is_delete             VARCHAR(1) DEFAULT 'N'
);

create table tbl_user_company_app(
                                     user_id_token         VARCHAR(36) primary key,
                                     in_active              VARCHAR(1) DEFAULT 'N',
                                     user_id               VARCHAR(50),
                                     app_code              VARCHAR(50),
                                     company_code          VARCHAR(50),
                                     branch_code           VARCHAR(50),
    -- Pattern Fields
                                     create_at             TIMESTAMPTZ DEFAULT now(),
                                     create_by             VARCHAR(50),
                                     update_by             VARCHAR(50),
                                     update_at             TIMESTAMPTZ,
                                     is_delete             VARCHAR(1) DEFAULT 'N',
                                     user_active_time      TIMESTAMPTZ,
                                     unique (app_code,company_code,user_id)
);
CREATE TABLE tbl_user_role
(
    role_code   VARCHAR(50)  REFERENCES tbl_role (role_code) ON DELETE CASCADE,
    user_id_token   VARCHAR(36) REFERENCES tbl_user_company_app(user_id_token) on delete cascade ,
    PRIMARY KEY (role_code, user_id_token),
    create_at TIMESTAMPTZ DEFAULT now(),
    create_by VARCHAR(50),
    update_by VARCHAR(50),
    update_at TIMESTAMPTZ,
    is_delete VARCHAR(1) DEFAULT 'N'
);

create INDEX idx_user_id_external_id on tbl_user (external_id);
create INDEX idx_user_company_status on tbl_user (status);
create INDEX idx_user_occupation on tbl_user (occupation);
create INDEX tbl_user_company_app_company_code_app_code on tbl_user_company_app (company_code, app_code);
create INDEX tbl_user_company_app_company_code_app_code_branch_code on tbl_user_company_app (company_code, app_code,branch_code);

GRANT ALL PRIVILEGES ON DATABASE home TO home;
GRANT ALL ON SCHEMA public TO home;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO home;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO home;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO home;