#------------------------------------------------------------
#        Script MySQL.
#------------------------------------------------------------


#------------------------------------------------------------
# Table: user
#------------------------------------------------------------

CREATE TABLE user(
        id       Int  Auto_increment  NOT NULL ,
        name      Varchar (50) NOT NULL ,
        admin    Bool NOT NULL ,
        password BINARY (32) NOT NULL ,
        salt     BINARY (32) NOT NULL
    ,CONSTRAINT user_PK PRIMARY KEY (id)
)ENGINE=InnoDB;


#------------------------------------------------------------
# Table: sonde
#------------------------------------------------------------

CREATE TABLE sonde(
        id        Int  Auto_increment  NOT NULL ,
        active    Bool NOT NULL ,
        latitude  Float NOT NULL ,
        longitude Float NOT NULL
    ,CONSTRAINT sonde_PK PRIMARY KEY (id)
)ENGINE=InnoDB;


#------------------------------------------------------------
# Table: relevé
#------------------------------------------------------------

CREATE TABLE releve(
        id          Int  Auto_increment  NOT NULL ,
        temperature Float NOT NULL ,
        date        Datetime NOT NULL default current_timestamp ,
        humidite    Decimal NOT NULL ,
        id_sonde    Int NOT NULL
    ,CONSTRAINT releve_PK PRIMARY KEY (id)

    ,CONSTRAINT releve_sonde_FK FOREIGN KEY (id_sonde) REFERENCES sonde(id)
)ENGINE=InnoDB;


#------------------------------------------------------------
# Table: posséde
#------------------------------------------------------------

CREATE TABLE possede(
        id_sonde Int NOT NULL ,
        id_user Int NOT NULL
    ,CONSTRAINT possede_PK PRIMARY KEY (id_sonde,id_user)

    ,CONSTRAINT possede_sonde_FK FOREIGN KEY (id_sonde) REFERENCES sonde(id)
    ,CONSTRAINT possede_user0_FK FOREIGN KEY (id_user) REFERENCES user(id)
)ENGINE=InnoDB;