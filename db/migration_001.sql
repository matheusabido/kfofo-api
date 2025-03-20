CREATE TABLE users (id SERIAL PRIMARY KEY, name VARCHAR(250) NOT NULL, email VARCHAR(250) NOT NULL UNIQUE, birth_date DATE NOT NULL, password TEXT NOT NULL);
CREATE TABLE restrictions (id SERIAL PRIMARY KEY, name VARCHAR(250) NOT NULL, description TEXT NOT NULL);
CREATE TABLE share_types (id SERIAL PRIMARY KEY, name VARCHAR(250) NOT NULL, description TEXT NOT NULL);
CREATE TABLE utensils (id SERIAL PRIMARY KEY, name VARCHAR(250) NOT NULL);
CREATE TABLE homes (id SERIAL PRIMARY KEY, user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE, address TEXT NOT NULL, city VARCHAR(250) NOT NULL, description TEXT NOT NULL, cost_day REAL NOT NULL, cost_week REAL, cost_month REAL, picture_path TEXT, restriction_id INT NOT NULL REFERENCES restrictions(id) ON DELETE RESTRICT, share_type_id INT NOT NULL REFERENCES share_types(id) ON DELETE RESTRICT);
CREATE TABLE bookings (id SERIAL PRIMARY KEY, user_id INT REFERENCES users(id) ON DELETE SET NULL, home_id INT REFERENCES homes(id) ON DELETE SET NULL, from_date DATE NOT NULL, to_date DATE NOT NULL, payment_type INT NOT NULL, COST_PER_CYCLE REAL NOT NUll);
CREATE TABLE home_utensils_pivot (home_id INT NOT NULL REFERENCES homes(id) ON DELETE CASCADE, utensil_id INT NOT NULL REFERENCES utensils(id) ON DELETE CASCADE, PRIMARY KEY(home_id, utensil_id));

INSERT INTO restrictions (name, description) VALUES ('Estudante', 'Essa casa admite somente estudantes. Comprovação é necessária.'), ('Trabalhador', 'Essa casa admite somente trabalhadores. Comprovação é necessária.'), ('Estudante/Trabalhador', 'Essa casa admite somente estudantes e trabalhadores. Comprovação é necessária.'), ('Nenhuma', 'Não há restrição nessa casa.');
INSERT INTO share_types (name, description) VALUES ('Família', 'Ideal para famílias. Ambiente privado, não há colegas de quarto.'), ('Compartilhado Misto', 'Essa casa é compartilhada. São aceitos homens e mulheres.'), ('Compartilhado Homem', 'Essa casa é compartilhada. São aceitos apenas homens.'), ('Compartilhado Mulher', 'Essa casa é compartilhada. São aceitas apenas mulheres');
INSERT INTO utensils (name) VALUES ('Cama'), ('Roupa de Cama'), ('Louça'), ('Climatização'), ('Lava-louças'), ('Máquina de lavar');
