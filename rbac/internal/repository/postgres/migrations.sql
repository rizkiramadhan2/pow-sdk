CREATE TABLE IF NOT EXISTS rbac_workspaces (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rbac_users (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    email TEXT,
    name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (workspace_id, id),

    CONSTRAINT fk_rbac_users_workspace
        FOREIGN KEY (workspace_id)
        REFERENCES rbac_workspaces(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rbac_roles (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (workspace_id, id),

    CONSTRAINT fk_rbac_roles_workspace
        FOREIGN KEY (workspace_id)
        REFERENCES rbac_workspaces(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rbac_policies (
    id TEXT NOT NULL,
    workspace_id TEXT NOT NULL,
    resource TEXT NOT NULL,
    action TEXT NOT NULL,
    effect TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (workspace_id, id),

    CONSTRAINT fk_rbac_policies_workspace
        FOREIGN KEY (workspace_id)
        REFERENCES rbac_workspaces(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_rbac_policies_effect
        CHECK (effect IN ('allow', 'deny'))
);

CREATE TABLE IF NOT EXISTS rbac_user_roles (
    workspace_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (workspace_id, user_id, role_id),

    CONSTRAINT fk_rbac_user_roles_user
        FOREIGN KEY (workspace_id, user_id)
        REFERENCES rbac_users(workspace_id, id)
        ON DELETE CASCADE,

    CONSTRAINT fk_rbac_user_roles_role
        FOREIGN KEY (workspace_id, role_id)
        REFERENCES rbac_roles(workspace_id, id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rbac_role_policies (
    workspace_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    policy_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (workspace_id, role_id, policy_id),

    CONSTRAINT fk_rbac_role_policies_role
        FOREIGN KEY (workspace_id, role_id)
        REFERENCES rbac_roles(workspace_id, id)
        ON DELETE CASCADE,

    CONSTRAINT fk_rbac_role_policies_policy
        FOREIGN KEY (workspace_id, policy_id)
        REFERENCES rbac_policies(workspace_id, id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_rbac_users_workspace_id
    ON rbac_users(workspace_id);

CREATE INDEX IF NOT EXISTS idx_rbac_roles_workspace_id
    ON rbac_roles(workspace_id);

CREATE INDEX IF NOT EXISTS idx_rbac_policies_workspace_id
    ON rbac_policies(workspace_id);

CREATE INDEX IF NOT EXISTS idx_rbac_user_roles_user
    ON rbac_user_roles(workspace_id, user_id);

CREATE INDEX IF NOT EXISTS idx_rbac_role_policies_role
    ON rbac_role_policies(workspace_id, role_id);