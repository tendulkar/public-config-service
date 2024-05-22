-- Types
INSERT INTO types (namespace, family, name, element_type, widget_type)
VALUES ('default', 'input', 'text', 'text', 'text_field'), 
    ('default', 'input', 'long_text', 'text', 'long_text_field'), 
    ('default', 'input', 'decimal', 'decimal', 'number_field'),
    ('default', 'input', 'number', 'number', 'number_field'),
    ('default', 'input', 'currency', 'currency', 'currency_field'), 
    ('default', 'input', 'quantity', 'quantity', 'quantity_field'),
    ('default', 'input', 'percentage', 'percentage', 'percentage_field'), 
    ('default', 'input', 'email', 'email', 'email_field'), 
    ('default', 'input', 'phone', 'phone', 'phone_field'),
    ('default', 'input', 'password', 'password', 'password_field'), 
    ('default', 'input', 'boolean', 'boolean', 'boolean_field'), 
    ('default', 'input', 'choice', 'choice', 'choice_field'),
    ('default', 'input', 'multi_choice', 'multi_choice', 'multi_choice_field'), 
    ('default', 'input', 'year', 'year', 'year_picker'),
    ('default', 'input', 'month', 'month', 'month_picker'), 
    ('default', 'input', 'date', 'date', 'date_picker'),
    ('default', 'input', 'day', 'day', 'day_pickeer'), 
    ('default', 'input', 'time', 'time', 'time_picker'),
    ('default', 'input', 'image', 'image', 'image_picker'),
    ('default', 'input', 'audio', 'audio', 'audio_picker'), 
    ('default', 'input', 'video', 'video', 'video_picker'),
    ('default', 'input', 'document', 'document', 'document_picker'), 
    ('default', 'input', 'contact', 'contact', 'contact_picker'),
    ('default', 'input', 'location_pin', 'location_pin', 'location_picker');

-- Validations
INSERT INTO validations (namespace, family, name, rule_name, validation_params)
VALUES ('default', 'text', 'min_length', 'min_length', '{}'),
    ('default', 'text', 'max_length', 'max_length', '{}'),
    ('default', 'text', 'required', 'required', '{}'),
    ('default', 'text', 'email', 'email', '{}'),
    ('default', 'text', 'numeric', 'numeric', '{}'),
    ('default', 'text', 'alpha_only', 'alpha_only', '{}'),
    ('default', 'text', 'alpha_numeric', 'alpha_numeric', '{}'),
    ('default', 'text', 'alpha_numeric_special', 'alpha_numeric_special', '{}'),
    ('default', 'text', 'match_pattern', 'match_pattern', '{}'),
    ('default', 'text', 'valid_url', 'valid_url', '{}'),
    ('default', 'text', 'valid_ip', 'valid_ip', '{}'),
    ('default', 'text', 'valid_credit_card', 'valid_credit_card', '{}'),
    ('default', 'text', 'exact_length', 'exact_length', '{}'),
    ('default', 'text', 'max_value', 'max_value', '{}'),
    ('default', 'text', 'min_value', 'min_value', '{}'),
    ('default', 'text', 'equals_nocase', 'equals_nocase', '{}'),
    ('default', 'text', 'equals', 'equals', '{}'),
    ('default', 'text', 'min_words_count', 'min_words_count', '{}'),
    ('default', 'text', 'max_words_count', 'max_words_count', '{}'),
    ('default', 'text', 'exact_words_count', 'exact_words_count', '{}');

-- Attributes
INSERT INTO attributes (namespace, family, name, label)
VALUES ('example', 'input', 'text', 'Example Text Input');

-- Forms
INSERT INTO forms (namespace, family, name, action_name, attributes)
VALUES ('example', 'form', 'example_form', 'submit', '{}');

-- Form Attributes
INSERT INTO form_attributes (form_id, attribute_id)
VALUES (1, 1);

-- Attribute Validations
INSERT INTO attribute_validations (attribute_id, validation_id)
VALUES (1, 1);