CREATE TABLE tbl_user
(
    user_id_token         VARCHAR(36) primary key , -- primary key
    user_id               VARCHAR(50),
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
    branch_code           VARCHAR(50),
    app_code              VARCHAR(50),
    company_code          VARCHAR(50),
    status                VARCHAR(50),
    account_name          VARCHAR(100),
    user_active_time      TIMESTAMPTZ,
    external_id           VARCHAR(100),
    user_details          JSON,
    in_active             BOOLEAN    DEFAULT FALSE,

    -- Pattern Fields
    create_at             TIMESTAMPTZ DEFAULT now(),
    create_by             VARCHAR(50),
    update_by             VARCHAR(50),
    update_at             TIMESTAMPTZ,
    is_delete             VARCHAR(1) DEFAULT 'N',
    unique (user_id, company_code, app_code)
);
ALTER TABLE tbl_user
ALTER
COLUMN is_active TYPE VARCHAR(1) USING (is_active::VARCHAR(1)),
    ALTER
COLUMN is_active SET DEFAULT 'N';


create INDEX idx_user_id_branch_company on tbl_user (user_id, branch_code, company_code);
create INDEX idx_user_company_branch on tbl_user (company_code, branch_code);
create INDEX idx_user_company on tbl_user (company_code);
create INDEX idx_user_company_app_code on tbl_user (company_code, app_code);



CREATE TABLE tbl_role
(
    role_code      VARCHAR(50) primary key, -- role code = prefix company + app + role code
    parent_role_code VARCHAR(50) REFERENCES tbl_role (role_code) ON DELETE SET NULL,
    role_name_th   VARCHAR(100),
    role_name_en   VARCHAR(100),
    role_desc_th   VARCHAR(255),
    role_desc_en   VARCHAR(255),
    create_at      TIMESTAMPTZ DEFAULT now(),
    create_by      VARCHAR(50),
    update_by      VARCHAR(50),
    update_at      TIMESTAMPTZ,
    is_delete      VARCHAR(1) DEFAULT 'N',
    unique (role_code)
);


CREATE TABLE tbl_user_role
(
    role_code   VARCHAR(50)  REFERENCES tbl_role (role_code) ON DELETE CASCADE,
    user_id_token   VARCHAR(36) REFERENCES tbl_user(user_id_token) on delete cascade ,
    PRIMARY KEY (role_code, user_id_token),
    create_at TIMESTAMPTZ DEFAULT now(),
    create_by VARCHAR(50),
    update_by VARCHAR(50),
    update_at TIMESTAMPTZ,
    is_delete VARCHAR(1) DEFAULT 'N'
);



CREATE TABLE tbl_object
(
    object_code VARCHAR(50) PRIMARY KEY,
    object_name VARCHAR(100),
    object_desc VARCHAR(100),
    create_at   TIMESTAMPTZ DEFAULT now(),
    create_by   VARCHAR(50),
    update_by   VARCHAR(50),
    update_at   TIMESTAMPTZ,
    is_delete   VARCHAR(1) DEFAULT 'N'
);


CREATE TABLE tbl_role_object
(
    id serir
    role_code   VARCHAR(50) REFERENCES tbl_role (role_code) ON DELETE CASCADE,
    object_code VARCHAR(50) REFERENCES tbl_object (object_code) ON DELETE CASCADE,
    PRIMARY KEY (role_code, object_code),
    create_at   TIMESTAMPTZ DEFAULT now(),
    create_by   VARCHAR(50),
    update_by   VARCHAR(50),
    update_at   TIMESTAMPTZ,
    is_delete   VARCHAR(1) DEFAULT 'N'
);

CREATE TABLE tbl_company
(
    company_code    VARCHAR(10) PRIMARY KEY,
    company_name_th VARCHAR(255),
    company_name_en VARCHAR(255),
    company_picture BYTEA,
    company_desc_th TEXT,
    company_desc_en TEXT,
    company_details jsonb,
    owner_id        VARCHAR(100),
    create_at       TIMESTAMPTZ DEFAULT now(),
    create_by       VARCHAR(50),
    update_by       VARCHAR(50),
    update_at       TIMESTAMPTZ,
    is_delete       VARCHAR(1) DEFAULT 'N'
);


CREATE TABLE tbl_branch
(
    branch_code    VARCHAR(15) PRIMARY KEY,
    branch_name_th VARCHAR(255),
    branch_name_en VARCHAR(255),
    branch_picture BYTEA,
    branch_desc_th TEXT,
    branch_desc_en TEXT,
    branch_details jsonb,
    owner_id        VARCHAR(100),
    create_at       TIMESTAMPTZ DEFAULT now(),
    create_by       VARCHAR(50),
    update_by       VARCHAR(50),
    update_at       TIMESTAMPTZ,
    is_delete       VARCHAR(1) DEFAULT 'N'
);

CREATE TABLE tbl_app
(
    app_code      VARCHAR(10),
    company_code  VARCHAR(10),
    app_name_th   VARCHAR(255),
    app_name_en   VARCHAR(255),
    app_picture   BYTEA,
    app_desc_th   TEXT,
    app_desc_en   TEXT,
    owner_id      VARCHAR(100),
    user_role     VARCHAR(50),
    user_approve  CHAR(1)    DEFAULT 'N',
    approved_type VARCHAR(50), -- default email
    create_at     TIMESTAMPTZ DEFAULT now(),
    create_by     VARCHAR(50),
    update_by     VARCHAR(50),
    update_at     TIMESTAMPTZ,
    is_delete     VARCHAR(1) DEFAULT 'N',
    PRIMARY KEY(company_code,app_code)
);



