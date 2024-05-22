CREATE TABLE types (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    family VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    element_type VARCHAR(255) NOT NULL,
    widget_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    deleted_at TIMESTAMP(3),
    version INTEGER DEFAULT 1 NOT NULL,
    UNIQUE (namespace, family, name)
);

CREATE TABLE validations (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    family VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    rule_name VARCHAR(255) NOT NULL,
    validation_params JSON,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    deleted_at TIMESTAMP(3),
    version INTEGER DEFAULT 1 NOT NULL,
    UNIQUE (namespace, family, name)
);

CREATE TABLE attributes (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    family VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    design_spec JSON,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    deleted_at TIMESTAMP(3),
    version INTEGER DEFAULT 1 NOT NULL,
    UNIQUE (namespace, family, name)
);

CREATE TABLE forms (
    id SERIAL PRIMARY KEY,
    namespace VARCHAR(255) NOT NULL,
    family VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    action_name VARCHAR(255) NOT NULL,
    attributes JSON NOT NULL,
    created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    deleted_at TIMESTAMP(3),
    version INTEGER DEFAULT 1 NOT NULL,
    UNIQUE (namespace, family, name)
);


-- Table: form_attributes (to link forms and attributes)
CREATE TABLE form_attributes (
    id SERIAL PRIMARY KEY,
    form_id INT NOT NULL,
    attribute_id INT NOT NULL,
    FOREIGN KEY (form_id) REFERENCES forms(id),
    FOREIGN KEY (attribute_id) REFERENCES attributes(id)
);

-- Table: attribute_validations (to link attributes and validations)
CREATE TABLE attribute_validations (
    id SERIAL PRIMARY KEY,
    attribute_id INT NOT NULL,
    validation_id INT NOT NULL,
    FOREIGN KEY (attribute_id) REFERENCES attributes(id),
    FOREIGN KEY (validation_id) REFERENCES validations(id)
);
