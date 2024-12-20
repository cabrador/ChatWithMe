package db

const schema = `
SET search_path TO public;

CREATE TABLE IF NOT EXISTS personality_traits (
    id SERIAL PRIMARY KEY,
	trait text  
);


CREATE TABLE IF NOT EXISTS personas (
    id SERIAL PRIMARY KEY,
    first_name text,
    last_name text
);

CREATE TABLE IF NOT EXISTS persona_personality_traits (
    id SERIAL PRIMARY KEY,
    persona_id INT NOT NULL,
    personality_trait_id INT NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(persona_id) REFERENCES personas(id),
    CONSTRAINT fk_personality_trait FOREIGN KEY(personality_trait_id) REFERENCES personality_traits(id)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
	username text
);

CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    author text
);
    
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    persona_id INT NOT NULL,
    content text,
    order_number INT,
    author_id INT,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
    CONSTRAINT fk_persona FOREIGN KEY(persona_id) REFERENCES personas(id),
    CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES authors(id)
);	
`
