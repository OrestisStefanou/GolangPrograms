CREATE TABLE Wallets(
        id VARCHAR(40) NOT NULL,
        publicKey VARCHAR(40) NOT NULL,
        privateKey VARCHAR(40) NOT NULL,
        PRIMARY KEY(id)
);  

CREATE TABLE Bitcoins(
	id VARCHAR(40) NOT NULL,
	firstOwner VARCHAR(40) NOT NULL,
	value INT NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY(firstOwner) REFERENCES Wallets(id)
);

CREATE TABLE Transactions(
	id INT UNSIGNED AUTO_INCREMENT NOT NULL,
        BitcoinID VARCHAR(40) NOT NULL,
        senderID VARCHAR(40) NOT NULL,
	receiverID VARCHAR(40) NOT NULL,
	value INT NOT NULL,
	creation_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        PRIMARY KEY(id),
        FOREIGN KEY(BitcoinID) REFERENCES Bitcoins(id),
	FOREIGN KEY(senderID) REFERENCES Wallets(id),
	FOREIGN KEY(receiverID) REFERENCES Wallets(id)
);
