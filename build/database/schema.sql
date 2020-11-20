USE nmap;

CREATE TABLE runs (
    id INT NOT NULL AUTO_INCREMENT,
    address VARCHAR(253) NOT NULL, -- Why 253? That should be the character limit, incl dots, of a hostname
    PRIMARY KEY (id, address)
)
;

CREATE TABLE run_results (
    id INT NOT NULL AUTO_INCREMENT,
    run INT NOT NULL,
    address VARCHAR(253) NOT NULL,
    ports TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY(run) REFERENCES runs(id)
)
;
