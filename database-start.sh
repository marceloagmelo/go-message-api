#!/usr/bin/env bash

# Tabela
echo "Criando a tabela mensagem..."
mysql -h localhost -u root -p -D gomessagedb << EOF
use gomessagedb;
CREATE TABLE mensagem (
id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
titulo VARCHAR(100), texto VARCHAR(255),
status INTEGER,
PRIMARY KEY (id)
);
EOF

