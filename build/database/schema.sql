USE nmap;

CREATE TABLE runs (
    id INT NOT NULL AUTO_INCREMENT,
    address VARCHAR(253) NOT NULL,
    ports TEXT NOT NULL,
    PRIMARY KEY (id)
)
;
